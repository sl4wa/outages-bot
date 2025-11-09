<?php

namespace App\Infrastructure\Telegram\Handlers;

use App\Application\Bot\Service\UnsubscribeUserService;
use SergiX44\Nutgram\Handlers\Type\Command;
use SergiX44\Nutgram\Nutgram;

class StopCommand extends Command
{
    protected string $command = 'stop';
    protected ?string $description = 'Відписатися від сповіщень';

    public function __construct(private readonly UnsubscribeUserService $unsubscribeUserService)
    {
        parent::__construct();
    }

    public function handle(Nutgram $bot): void
    {
        $bot->sendMessage($this->unsubscribeUserService->handle($bot->chatId()));
    }
}
