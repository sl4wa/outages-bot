<?php

namespace App\Infrastructure\Telegram\Handlers;

use App\Infrastructure\Repository\FileUserRepository;
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

    private FileUserRepository $userRepository;
    private FileStreetRepository $streetRepository;

    // Data persists between steps
    public int $selectedStreetId = 0;
    public string $selectedStreetName = '';

    public function __construct(
        FileUserRepository $userRepository,
        FileStreetRepository $streetRepository
    ) {
        $this->userRepository = $userRepository;
        $this->streetRepository = $streetRepository;
    }

    public function askStreet(Nutgram $bot)
    {
        $user = $this->userRepository->find($bot->chatId());
        if ($user) {
            $bot->sendMessage(
                "Ваша поточна підписка:\nВулиця: {$user->streetName}\nБудинок: {$user->building}\n\n"
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

        $user = new User(
            id: $bot->chatId(),
            streetId: $this->selectedStreetId,
            streetName: $this->selectedStreetName,
            building: $building,
            startDate: null,
            endDate: null,
            comment: ''
        );
        $this->userRepository->save($user);

        $bot->sendMessage(
            "Ви підписалися на сповіщення про відключення електроенергії для вулиці {$user->streetName}, будинок {$building}.",
            reply_markup: ReplyKeyboardRemove::make(true)
        );
        $this->end();
    }
}
