<?php

declare(strict_types=1);

namespace App\Application\Bot\Service;

use App\Application\Bot\DTO\SearchStreetResultDTO;
use App\Application\Bot\Query\FilterStreetQueryHandler;

final readonly class SearchStreetService
{
    public function __construct(
        private FilterStreetQueryHandler $filterStreetsQueryHandler
    ) {
    }

    public function handle(string $query): SearchStreetResultDTO
    {
        $query = trim($query);

        if ($query === '') {
            return new SearchStreetResultDTO(
                message: 'Введіть назву вулиці.',
                shouldContinue: false
            );
        }

        $streets = $this->filterStreetsQueryHandler->handle($query);

        if ($streets === []) {
            return new SearchStreetResultDTO(
                message: 'Вулицю не знайдено. Спробуйте ще раз.',
                shouldContinue: false
            );
        }

        if (count($streets) === 1) {
            $street = $streets[0];

            return new SearchStreetResultDTO(
                message: "Ви обрали вулицю: {$street->name}\nБудь ласка, введіть номер будинку:",
                selectedStreetId: $street->id,
                selectedStreetName: $street->name,
                shouldContinue: true
            );
        }

        return new SearchStreetResultDTO(
            message: 'Будь ласка, оберіть вулицю:',
            streetOptions: $streets,
            shouldContinue: false
        );
    }
}
