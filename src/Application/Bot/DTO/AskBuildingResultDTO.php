<?php

declare(strict_types=1);

namespace App\Application\Bot\DTO;

final readonly class AskBuildingResultDTO
{
    public function __construct(
        public string $message,
        public bool $isSuccess
    ) {
    }
}
