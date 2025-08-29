<?php

namespace App\Application\DTO;

class UserSubscriptionDTO
{
    public function __construct(
        public readonly int $chatId,
        public readonly int $streetId,
        public readonly string $streetName,
        public readonly string $building,
    ) {}
}

