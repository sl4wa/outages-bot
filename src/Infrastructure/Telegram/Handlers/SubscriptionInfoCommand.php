<?php

namespace App\Infrastructure\Telegram\Handlers;

use App\Infrastructure\Repository\FileUserRepository;
use SergiX44\Nutgram\Handlers\Type\Command;
use SergiX44\Nutgram\Nutgram;

class SubscriptionInfoCommand extends Command
{
    protected string $command = 'subscription';
    protected ?string $description = 'Показати поточну підписку';

    private FileUserRepository $userRepository;

    public function __construct(FileUserRepository $userRepository)
    {
        $this->userRepository = $userRepository;
        parent::__construct();
    }

    public function handle(Nutgram $bot): void
    {
        $user = $this->userRepository->find($bot->chatId());
        if ($user) {
            $bot->sendMessage("Ваша поточна підписка:\nВулиця: {$user->streetName}\nБудинок: {$user->building}");
        } else {
            $bot->sendMessage("Ви не маєте активної підписки.");
        }
    }
}
