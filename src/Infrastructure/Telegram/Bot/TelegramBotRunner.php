<?php

declare(strict_types=1);

namespace App\Infrastructure\Telegram\Bot;

use App\Application\Interface\BotRunnerInterface;
use App\Infrastructure\Telegram\Handlers\StartHandler;
use App\Infrastructure\Telegram\Handlers\StopHandler;
use App\Infrastructure\Telegram\Handlers\SubscriptionHandler;
use SergiX44\Nutgram\Conversations\Conversation;
use SergiX44\Nutgram\Nutgram;
use Symfony\Component\DependencyInjection\ContainerInterface;

final class TelegramBotRunner implements BotRunnerInterface
{
    private Nutgram $bot;

    public function __construct(
        Nutgram $bot,
        ContainerInterface $container
    ) {
        $this->bot = $bot;
        $this->bot->getContainer()->delegate($container);
    }

    public function run(): void
    {
        Conversation::refreshOnDeserialize();

        $this->bot->onCommand('start', StartHandler::class);
        $this->bot->onCommand('stop', StopHandler::class);
        $this->bot->onCommand('subscription', SubscriptionHandler::class);

        $this->bot->run();
    }
}
