<?php

namespace App\Application\Bot\Service\Subscription;

use App\Application\Bot\DTO\SelectStreetResultDTO;
use App\Application\Interface\Repository\StreetRepositoryInterface;

readonly class SelectStreetService
{
    public function __construct(
        private StreetRepositoryInterface $streetRepository
    ) {
    }

    public function handle(string $query): SelectStreetResultDTO
    {
        $query = trim($query);

        // Validate empty input
        if ($query === '') {
            return new SelectStreetResultDTO(
                message: 'Введіть назву вулиці.',
                shouldContinue: false
            );
        }

        // Search for streets
        $filtered = $this->streetRepository->filter($query);

        // No streets found
        if (count($filtered) === 0) {
            return new SelectStreetResultDTO(
                message: 'Вулицю не знайдено. Спробуйте ще раз.',
                shouldContinue: false
            );
        }

        // Check for exact match
        $exact = $this->streetRepository->findByName($query);
        if ($exact) {
            return new SelectStreetResultDTO(
                message: "Ви обрали вулицю: {$exact['name']}\nБудь ласка, введіть номер будинку (наприклад: 13 або 13-А):",
                selectedStreetId: $exact['id'],
                selectedStreetName: $exact['name'],
                shouldContinue: true
            );
        }

        // Multiple matches - return options for keyboard
        return new SelectStreetResultDTO(
            message: 'Будь ласка, оберіть вулицю:',
            streetOptions: $filtered,
            shouldContinue: false
        );
    }
}
