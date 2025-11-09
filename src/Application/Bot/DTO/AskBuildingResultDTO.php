<?php

namespace App\Application\Bot\DTO;

readonly class AskBuildingResultDTO
{
    public function __construct(
        public string $message,
        public bool $isSuccess
    ) {
    }
}
