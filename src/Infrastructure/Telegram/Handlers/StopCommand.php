<?php

declare(strict_types=1);

namespace App\Infrastructure\Telegram\Handlers;

use App\Application\Bot\Command\UnsubscribeUserCommandHandler;
use SergiX44\Nutgram\Handlers\Type\Command;
use SergiX44\Nutgram\Nutgram;

final class StopCommand extends Command
{
    protected string $command = 'stop';

    protected ?string $description = 'Відписатися від сповіщень';

    public function __construct(private readonly UnsubscribeUserCommandHandler $unsubscribeUserCommandHandler)
    {
        parent::__construct();
    }

    public function handle(Nutgram $bot): void
    {
        $removed = $this->unsubscribeUserCommandHandler->handle($bot->chatId());

        $message = $removed
            ? 'Ви успішно відписалися від сповіщень про відключення електроенергії.'
            : 'Ви не маєте активної підписки.';

        $bot->sendMessage($message);
    }
}
