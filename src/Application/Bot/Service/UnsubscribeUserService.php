<?php

namespace App\Application\Bot\Service;

use App\Application\Bot\Command\UnsubscribeUserCommandHandler;

readonly class UnsubscribeUserService
{
    public function __construct(
        private UnsubscribeUserCommandHandler $unsubscribeUserCommandHandler
    ) {
    }

    public function handle(int $chatId): string
    {
        $removed = $this->unsubscribeUserCommandHandler->handle($chatId);

        if ($removed) {
            return 'Ви успішно відписалися від сповіщень про відключення електроенергії.';
        }

        return 'Ви не маєте активної підписки.';
    }
}
