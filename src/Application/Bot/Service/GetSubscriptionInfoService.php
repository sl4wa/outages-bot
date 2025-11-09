<?php

namespace App\Application\Bot\Service;

use App\Application\Bot\Query\GetUserSubscriptionQueryHandler;

readonly class GetSubscriptionInfoService
{
    public function __construct(
        private GetUserSubscriptionQueryHandler $getUserSubscriptionQueryHandler
    ) {
    }

    public function handle(int $chatId): string
    {
        $subscription = $this->getUserSubscriptionQueryHandler->handle($chatId);

        if ($subscription) {
            return "Ваша поточна підписка:\nВулиця: {$subscription->streetName}\nБудинок: {$subscription->building}";
        }

        return "Ви не маєте активної підписки.";
    }
}
