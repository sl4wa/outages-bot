<?php

declare(strict_types=1);

namespace App\Application\Bot\DTO;

final readonly class UserSubscriptionDTO
{
    public function __construct(
        public int $chatId,
        public int $streetId,
        public string $streetName,
        public string $building,
    ) {
    }
}
