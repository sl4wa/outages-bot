<?php

namespace App\Application\DTO;

class OutageDTO
{
    public function __construct(
        public readonly \DateTimeImmutable $start,
        public readonly \DateTimeImmutable $end,
        public readonly string $city,
        public readonly int $streetId,
        public readonly string $streetName,
        /** @var string[] */
        public readonly array $buildingNames,
        public readonly string $comment,
    ) {}
}
