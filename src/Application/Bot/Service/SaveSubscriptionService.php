<?php

declare(strict_types=1);

namespace App\Application\Bot\Service;

use App\Application\Bot\Command\CreateOrUpdateUserSubscriptionCommandHandler;
use App\Application\Bot\DTO\SaveSubscriptionResultDTO;
use App\Domain\Exception\InvalidBuildingFormatException;

final readonly class SaveSubscriptionService
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
    ): SaveSubscriptionResultDTO {
        // Validate state
        if (!$selectedStreetId || !$selectedStreetName) {
            return new SaveSubscriptionResultDTO(
                message: 'Підписка не завершена. Будь ласка, почніть знову.',
                isSuccess: false
            );
        }

        try {
            $result = $this->createOrUpdateUserSubscriptionCommandHandler->handle(
                chatId: $chatId,
                streetId: $selectedStreetId,
                streetName: $selectedStreetName,
                building: $building,
            );

            return new SaveSubscriptionResultDTO(
                message: "Ви підписалися на сповіщення про відключення електроенергії для вулиці {$result->streetName}, будинок {$result->building}.",
                isSuccess: true
            );
        } catch (InvalidBuildingFormatException) {
            return new SaveSubscriptionResultDTO(
                message: 'Невірний формат номера будинку. Приклад: 13 або 13-А',
                isSuccess: false
            );
        }
    }
}
