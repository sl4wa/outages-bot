<?php

namespace App\Infrastructure\Telegram\Handlers;

use App\Infrastructure\Repository\FileUserRepository;
use SergiX44\Nutgram\Handlers\Type\Command;
use SergiX44\Nutgram\Nutgram;

class StopCommand extends Command
{
    protected string $command = 'stop';
    protected ?string $description = 'Відписатися від сповіщень';

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
            $this->userRepository->remove($bot->chatId());
            $bot->sendMessage('Ви успішно відписалися від сповіщень про відключення електроенергії.');
        } else {
            $bot->sendMessage('Ви не маєте активної підписки.');
        }
    }
}
