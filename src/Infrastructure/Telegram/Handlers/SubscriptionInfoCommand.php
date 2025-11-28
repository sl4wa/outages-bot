<?php

declare(strict_types=1);

namespace App\Infrastructure\Telegram\Handlers;

use App\Application\Bot\Query\GetUserSubscriptionQueryHandler;
use SergiX44\Nutgram\Handlers\Type\Command;
use SergiX44\Nutgram\Nutgram;

final class SubscriptionInfoCommand extends Command
{
    protected string $command = 'subscription';

    protected ?string $description = 'Показати поточну підписку';

    public function __construct(private readonly GetUserSubscriptionQueryHandler $getUserSubscriptionQueryHandler)
    {
        parent::__construct();
    }

    public function handle(Nutgram $bot): void
    {
        $subscription = $this->getUserSubscriptionQueryHandler->handle($bot->chatId());

        $message = $subscription
            ? "Ваша поточна підписка:\nВулиця: {$subscription->streetName}\nБудинок: {$subscription->building}"
            : 'Ви не маєте активної підписки.';

        $bot->sendMessage($message);
    }
}
