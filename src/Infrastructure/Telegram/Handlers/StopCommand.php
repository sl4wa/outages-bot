<?php

declare(strict_types=1);

namespace App\Infrastructure\Telegram\Handlers;

use App\Application\Interface\Repository\UserRepositoryInterface;
use SergiX44\Nutgram\Handlers\Type\Command;
use SergiX44\Nutgram\Nutgram;

final class StopCommand extends Command
{
    protected string $command = 'stop';

    protected ?string $description = 'Відписатися від сповіщень';

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
