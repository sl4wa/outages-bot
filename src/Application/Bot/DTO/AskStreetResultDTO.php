<?php

namespace App\Application\Bot\DTO;

readonly class AskStreetResultDTO
{
    public function __construct(
        public string $message
    ) {
    }
}
