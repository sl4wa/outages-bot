<?php

namespace App\Application\Bot\Service\Subscription;

use App\Application\Bot\Command\CreateOrUpdateUserSubscriptionCommandHandler;
use App\Application\Bot\DTO\AskBuildingResultDTO;

readonly class AskBuildingService
{
    public function __construct(
        private CreateOrUpdateUserSubscriptionCommandHandler $createOrUpdateUserSubscriptionCommandHandler
    ) {
    }

    public function handle(
        int $chatId,
        ?int $selectedStreetId,
        ?string $selectedStreetName,
        string $building
    ): AskBuildingResultDTO {
        // Validate state
        if (!$selectedStreetId || !$selectedStreetName) {
            return new AskBuildingResultDTO(
                message: 'Підписка не завершена. Будь ласка, почніть знову.',
                isSuccess: false
            );
        }

        // Validate building input
        $building = trim($building);
        if ($building === '') {
            return new AskBuildingResultDTO(
                message: 'Введіть номер будинку.',
                isSuccess: false
            );
        }

        // Create/update subscription
        $result = $this->createOrUpdateUserSubscriptionCommandHandler->handle(
            chatId: $chatId,
            streetId: $selectedStreetId,
            streetName: $selectedStreetName,
            building: $building,
        );

        return new AskBuildingResultDTO(
            message: "Ви підписалися на сповіщення про відключення електроенергії для вулиці {$result->streetName}, будинок {$result->building}.",
            isSuccess: true
        );
    }
}
