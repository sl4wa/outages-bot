<?php

declare(strict_types=1);

namespace App\Infrastructure\Telegram\Handlers;

use App\Application\Bot\Service\SaveSubscriptionService;
use App\Application\Bot\Service\SearchStreetService;
use App\Application\Bot\Service\ShowSubscriptionService;
use SergiX44\Nutgram\Conversations\Conversation;
use SergiX44\Nutgram\Nutgram;
use SergiX44\Nutgram\Telegram\Types\Keyboard\KeyboardButton;
use SergiX44\Nutgram\Telegram\Types\Keyboard\ReplyKeyboardMarkup;
use SergiX44\Nutgram\Telegram\Types\Keyboard\ReplyKeyboardRemove;
use Symfony\Component\DependencyInjection\Attribute\Autoconfigure;

#[Autoconfigure(public: true, shared: false)]
final class StartHandler extends Conversation
{
    // Data persists between steps
    public int $selectedStreetId = 0;

    public string $selectedStreetName = '';

    public function __construct(
        private readonly ShowSubscriptionService $showSubscriptionService,
        private readonly SearchStreetService $searchStreetService,
        private readonly SaveSubscriptionService $saveSubscriptionService
    ) {
    }

    public function start(Nutgram $bot): void
    {
        $chatId = $bot->chatId();

        if ($chatId === null) {
            $this->end();

            return;
        }

        $message = $this->showSubscriptionService->handle($chatId);
        $bot->sendMessage($message);
        $this->next('searchStreet');
    }

    public function searchStreet(Nutgram $bot): void
    {
        $query = $bot->message()->text ?? '';
        $result = $this->searchStreetService->handle($query);

        // Handle exact match - move to next step
        if ($result->hasExactMatch() && $result->selectedStreetId !== null && $result->selectedStreetName !== null) {
            $this->selectedStreetId = $result->selectedStreetId;
            $this->selectedStreetName = $result->selectedStreetName;
            $bot->sendMessage($result->message, reply_markup: ReplyKeyboardRemove::make(true));
            $this->next('saveSubscription');

            return;
        }

        // Handle multiple street options - build keyboard
        if ($result->hasMultipleOptions() && $result->streetOptions !== null) {
            $replyMarkup = ReplyKeyboardMarkup::make(
                resize_keyboard: true,
                one_time_keyboard: true
            );

            foreach ($result->streetOptions as $street) {
                $replyMarkup->addRow(KeyboardButton::make($street->name));
            }
            $bot->sendMessage($result->message, reply_markup: $replyMarkup);
        } else {
            // Error cases (empty input, not found)
            $bot->sendMessage($result->message);
        }

        $this->next('searchStreet');
    }

    public function saveSubscription(Nutgram $bot): void
    {
        $chatId = $bot->chatId();

        if ($chatId === null) {
            $this->end();

            return;
        }

        $building = $bot->message()->text ?? '';

        $result = $this->saveSubscriptionService->handle(
            chatId: $chatId,
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
                $this->next('saveSubscription');
            }
        }
    }
}
