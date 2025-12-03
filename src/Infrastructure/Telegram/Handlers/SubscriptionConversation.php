<?php

declare(strict_types=1);

namespace App\Infrastructure\Telegram\Handlers;

use App\Application\Bot\Service\AskBuildingService;
use App\Application\Bot\Service\AskStreetService;
use App\Application\Bot\Service\SelectStreetService;
use SergiX44\Nutgram\Conversations\Conversation;
use SergiX44\Nutgram\Nutgram;
use SergiX44\Nutgram\Telegram\Types\Keyboard\KeyboardButton;
use SergiX44\Nutgram\Telegram\Types\Keyboard\ReplyKeyboardMarkup;
use SergiX44\Nutgram\Telegram\Types\Keyboard\ReplyKeyboardRemove;

final class SubscriptionConversation extends Conversation
{
    protected ?string $step = 'askStreet';

    // Data persists between steps
    public int $selectedStreetId = 0;

    public string $selectedStreetName = '';

    public function __construct(
        private readonly AskStreetService $askStreetService,
        private readonly SelectStreetService $selectStreetService,
        private readonly AskBuildingService $askBuildingService
    ) {
    }

    public function askStreet(Nutgram $bot): void
    {
        $result = $this->askStreetService->handle($bot->chatId());
        $bot->sendMessage($result->message);
        $this->next('selectStreet');
    }

    public function selectStreet(Nutgram $bot): void
    {
        $query = $bot->message()->text ?? '';
        $result = $this->selectStreetService->handle($query);

        // Handle exact match - move to next step
        if ($result->hasExactMatch()) {
            $this->selectedStreetId = $result->selectedStreetId;
            $this->selectedStreetName = $result->selectedStreetName;
            $bot->sendMessage($result->message, reply_markup: ReplyKeyboardRemove::make(true));
            $this->next('askBuilding');

            return;
        }

        // Handle multiple street options - build keyboard
        if ($result->hasMultipleOptions()) {
            $replyMarkup = ReplyKeyboardMarkup::make(
                resize_keyboard: true,
                one_time_keyboard: true
            );

            foreach ($result->streetOptions as $street) {
                $replyMarkup->addRow(KeyboardButton::make($street['name']));
            }
            $bot->sendMessage($result->message, reply_markup: $replyMarkup);
        } else {
            // Error cases (empty input, not found)
            $bot->sendMessage($result->message);
        }

        $this->next('selectStreet');
    }

    public function askBuilding(Nutgram $bot): void
    {
        $building = $bot->message()->text ?? '';

        $result = $this->askBuildingService->handle(
            chatId: $bot->chatId(),
            selectedStreetId: $this->selectedStreetId,
            selectedStreetName: $this->selectedStreetName,
            building: $building
        );

        if ($result->isSuccess) {
            $bot->sendMessage($result->message, reply_markup: ReplyKeyboardRemove::make(true));
            $this->end();
        } else {
            $bot->sendMessage($result->message);

            // If state validation failed, end conversation; otherwise retry
            if (!$this->selectedStreetId || !$this->selectedStreetName) {
                $this->end();
            } else {
                $this->next('askBuilding');
            }
        }
    }
}
