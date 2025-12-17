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
        int $streetId,
        string $streetName,
        string $building
    ): SaveSubscriptionResultDTO {
        try {
            $result = $this->createOrUpdateUserSubscriptionCommandHandler->handle(
                chatId: $chatId,
                streetId: $streetId,
                streetName: $streetName,
                building: $building,
            );

            return new SaveSubscriptionResultDTO(
                message: "Ви підписалися на сповіщення про відключення електроенергії для вулиці {$result->streetName}, будинок {$result->building}.",
                success: true,
            );
        } catch (InvalidBuildingFormatException $e) {
            return new SaveSubscriptionResultDTO(
                message: $e->getMessage(),
                success: false,
            );
        }
    }
}
