<?php

namespace App\Application\Bot\Service\Subscription;

use App\Application\Bot\DTO\AskStreetResultDTO;
use App\Application\Bot\Query\GetUserSubscriptionQueryHandler;

readonly class AskStreetService
{
    public function __construct(
        private GetUserSubscriptionQueryHandler $getUserSubscriptionQueryHandler
    ) {
    }

    public function handle(int $chatId): AskStreetResultDTO
    {
        $subscription = $this->getUserSubscriptionQueryHandler->handle($chatId);

        if ($subscription) {
            $message = "Ваша поточна підписка:\nВулиця: {$subscription->streetName}\nБудинок: {$subscription->building}\n\n"
                . "Будь ласка, оберіть нову вулицю для оновлення підписки або введіть назву вулиці:";
        } else {
            $message = "Будь ласка, введіть назву вулиці:";
        }

        return new AskStreetResultDTO($message);
    }
}
