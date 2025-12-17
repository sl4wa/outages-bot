<?php

declare(strict_types=1);

namespace App\Application\Bot\Service;

use App\Application\Bot\Query\GetUserSubscriptionQueryHandler;
use Throwable;

final readonly class ShowSubscriptionService
{
    public function __construct(
        private GetUserSubscriptionQueryHandler $getUserSubscriptionQueryHandler
    ) {
    }

    public function handle(int $chatId): string
    {
        try {
            $subscription = $this->getUserSubscriptionQueryHandler->handle($chatId);
        } catch (Throwable) {
            // If there's any error loading existing subscription (e.g., corrupted data),
            // ignore it and treat as new subscription
            $subscription = null;
        }

        if ($subscription) {
            return "Ваша поточна підписка:\nВулиця: {$subscription->streetName}\nБудинок: {$subscription->building}\n\n"
                . 'Будь ласка, введіть нову назву вулиці для оновлення підписки:';
        }

        return 'Будь ласка, введіть назву вулиці:';
    }
}
