<?php

declare(strict_types=1);

namespace App\Infrastructure\Telegram\Handlers;

use App\Domain\Repository\UserRepositoryInterface;
use SergiX44\Nutgram\Handlers\Type\Command;
use SergiX44\Nutgram\Nutgram;
use Symfony\Component\DependencyInjection\Attribute\Autoconfigure;

#[Autoconfigure(public: true)]
final class StopHandler extends Command
{
    public function __construct(private readonly UserRepositoryInterface $userRepository)
    {
        parent::__construct();
    }

    public function handle(Nutgram $bot): void
    {
        $chatId = $bot->chatId();

        if ($chatId === null) {
            return;
        }

        $removed = $this->userRepository->remove($chatId);

        $message = $removed
            ? 'Ви успішно відписалися від сповіщень про відключення електроенергії.'
            : 'Ви не маєте активної підписки.';

        $bot->sendMessage($message);
    }
}
