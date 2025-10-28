<?php

namespace App\Application\Bot\DTO;

readonly class UserSubscriptionDTO
{
    public function __construct(
        public int $chatId,
        public int $streetId,
        public string $streetName,
        public string $building,
    ) {}
}

