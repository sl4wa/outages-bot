<?php

declare(strict_types=1);

namespace App\Infrastructure\Telegram\Handlers;

use App\Application\Bot\Query\GetUserSubscriptionQueryHandler;
use SergiX44\Nutgram\Handlers\Type\Command;
use SergiX44\Nutgram\Nutgram;
use Symfony\Component\DependencyInjection\Attribute\Autoconfigure;

#[Autoconfigure(public: true)]
final class SubscriptionHandler extends Command
{
    public function __construct(private readonly GetUserSubscriptionQueryHandler $getUserSubscriptionQueryHandler)
    {
        parent::__construct();
    }

    public function handle(Nutgram $bot): void
    {
        $chatId = $bot->chatId();

        if ($chatId === null) {
            return;
        }

        $subscription = $this->getUserSubscriptionQueryHandler->handle($chatId);

        $message = $subscription
            ? "Ваша поточна підписка:\nВулиця: {$subscription->streetName}\nБудинок: {$subscription->building}"
            : 'Ви не маєте активної підписки.';

        $bot->sendMessage($message);
    }
}
