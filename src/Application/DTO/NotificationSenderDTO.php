<?php

namespace App\Application\DTO;

readonly class NotificationSenderDTO
{
    /**
     * @param array<int, string> $buildings
     */
    public function __construct(
        public int $userId,
        public string $city,
        public string $streetName,
        public array $buildings,
        public \DateTimeImmutable $start,
        public \DateTimeImmutable $end,
        public string $comment,
    ) {}
}
