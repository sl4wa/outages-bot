<?php

declare(strict_types=1);

namespace App\Application\Bot\DTO;

final readonly class AskStreetResultDTO
{
    public function __construct(
        public string $message
    ) {
    }
}
