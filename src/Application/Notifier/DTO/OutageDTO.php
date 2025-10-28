<?php

namespace App\Application\Notifier\DTO;

readonly class OutageDTO
{
    public function __construct(
        public \DateTimeImmutable $start,
        public \DateTimeImmutable $end,
        public string $city,
        public int $streetId,
        public string $streetName,
        /** @var string[] */
        public array $buildings,
        public string $comment,
    ) {}
}
