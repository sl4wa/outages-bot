<?php

declare(strict_types=1);

namespace App\Application\Bot\DTO;

final readonly class SaveSubscriptionResultDTO
{
    public function __construct(
        public string $message,
        public bool $success = true,
    ) {
    }
}
