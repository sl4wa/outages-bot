<?php

declare(strict_types=1);

namespace App\Application\Bot\Service\Subscription;

use App\Application\Bot\DTO\AskStreetResultDTO;
use App\Application\Bot\Query\GetUserSubscriptionQueryHandler;

final readonly class AskStreetService
{
    public function __construct(
        private GetUserSubscriptionQueryHandler $getUserSubscriptionQueryHandler
    ) {
    }

    public function handle(int $chatId): AskStreetResultDTO
    {
        try {
            $subscription = $this->getUserSubscriptionQueryHandler->handle($chatId);
        } catch (\Throwable) {
            // If there's any error loading existing subscription (e.g., corrupted data),
            // ignore it and treat as new subscription
            $subscription = null;
        }

        if ($subscription) {
            $message = "Ваша поточна підписка:\nВулиця: {$subscription->streetName}\nБудинок: {$subscription->building}\n\n"
                . 'Будь ласка, оберіть нову вулицю для оновлення підписки або введіть назву вулиці:';
        } else {
            $message = 'Будь ласка, введіть назву вулиці:';
        }

        return new AskStreetResultDTO($message);
    }
}
