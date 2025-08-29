<?php

namespace App\Infrastructure\Telegram\Handlers;

use App\Application\Service\UnsubscribeService;
use SergiX44\Nutgram\Handlers\Type\Command;
use SergiX44\Nutgram\Nutgram;

class StopCommand extends Command
{
    protected string $command = 'stop';
    protected ?string $description = 'Відписатися від сповіщень';

    public function __construct(private readonly UnsubscribeService $unsubscribeService)
    {
        parent::__construct();
    }

    public function handle(Nutgram $bot): void
    {
        $removed = $this->unsubscribeService->unsubscribe($bot->chatId());
        if ($removed) {
            $bot->sendMessage('Ви успішно відписалися від сповіщень про відключення електроенергії.');
        } else {
            $bot->sendMessage('Ви не маєте активної підписки.');
        }
    }
}
