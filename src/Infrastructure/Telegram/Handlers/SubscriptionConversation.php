<?php

namespace App\Infrastructure\Telegram\Handlers;

use App\Application\Bot\Query\GetUserSubscriptionQueryHandler;
use App\Application\Bot\Command\CreateOrUpdateUserSubscriptionCommandHandler;
use App\Infrastructure\Repository\FileStreetRepository;
use App\Domain\Entity\User;
use SergiX44\Nutgram\Conversations\Conversation;
use SergiX44\Nutgram\Nutgram;
use SergiX44\Nutgram\Telegram\Types\Keyboard\ReplyKeyboardMarkup;
use SergiX44\Nutgram\Telegram\Types\Keyboard\ReplyKeyboardRemove;
use SergiX44\Nutgram\Telegram\Types\Keyboard\KeyboardButton;

class SubscriptionConversation extends Conversation
{
    protected ?string $step = 'askStreet';

    private FileStreetRepository $streetRepository;

    public function __construct(
        private readonly GetUserSubscriptionQueryHandler $getUserSubscriptionQueryHandler,
        private readonly CreateOrUpdateUserSubscriptionCommandHandler $createOrUpdateUserSubscriptionCommandHandler,
        FileStreetRepository $streetRepository
    ) {
        $this->streetRepository = $streetRepository;
    }

    // Data persists between steps
    public int $selectedStreetId = 0;
    public string $selectedStreetName = '';

    public function askStreet(Nutgram $bot)
    {
        $sub = $this->getUserSubscriptionQueryHandler->handle($bot->chatId());
        if ($sub) {
            $bot->sendMessage(
                "Ваша поточна підписка:\nВулиця: {$sub->streetName}\nБудинок: {$sub->building}\n\n"
                ."Будь ласка, оберіть нову вулицю для оновлення підписки або введіть назву вулиці:"
            );
        } else {
            $bot->sendMessage("Будь ласка, введіть назву вулиці:");
        }
        $this->next('selectStreet');
    }

    public function selectStreet(Nutgram $bot)
    {
        $query = trim($bot->message()->text ?? '');
        if ($query === '') {
            $bot->sendMessage('Введіть назву вулиці.');
            $this->next('selectStreet');
            return;
        }

        $filtered = $this->streetRepository->filter($query);
        if (count($filtered) === 0) {
            $bot->sendMessage('Вулицю не знайдено. Спробуйте ще раз.');
            $this->next('selectStreet');
            return;
        }

        // Exact match?
        $exact = $this->streetRepository->findByName($query);
        if ($exact) {
            $this->selectedStreetId = $exact['id'];
            $this->selectedStreetName = $exact['name'];
            $bot->sendMessage(
                "Ви обрали вулицю: {$exact['name']}\nБудь ласка, введіть номер будинку:",
                reply_markup: ReplyKeyboardRemove::make(true)
            );
            $this->next('askBuilding');
            return;
        }

        $replyMarkup = ReplyKeyboardMarkup::make(
            resize_keyboard: true,
            one_time_keyboard: true
        );
        foreach ($filtered as $st) {
            $replyMarkup->addRow(KeyboardButton::make($st['name']));
        }

        $bot->sendMessage(
            'Будь ласка, оберіть вулицю:',
            reply_markup: $replyMarkup
        );
        $this->next('selectStreet');
    }

    public function askBuilding(Nutgram $bot)
    {
        $building = trim($bot->message()->text ?? '');

        if (!$this->selectedStreetId || !$this->selectedStreetName) {
            $bot->sendMessage('Підписка не завершена. Будь ласка, почніть знову.');
            $this->end();
            return;
        }
        if ($building === '') {
            $bot->sendMessage('Введіть номер будинку.');
            $this->next('askBuilding');
            return;
        }

        $result = $this->createOrUpdateUserSubscriptionCommandHandler->handle(
            chatId: $bot->chatId(),
            streetId: $this->selectedStreetId,
            streetName: $this->selectedStreetName,
            building: $building,
        );

        $bot->sendMessage(
            "Ви підписалися на сповіщення про відключення електроенергії для вулиці {$result->streetName}, будинок {$result->building}.",
            reply_markup: ReplyKeyboardRemove::make(true)
        );
        $this->end();
    }
}
