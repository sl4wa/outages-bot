<?php

declare(strict_types=1);

namespace App\Application\Notifier\DTO;

use DateTimeImmutable;

final readonly class OutageDTO
{
    public function __construct(
        public int $id,
        public DateTimeImmutable $start,
        public DateTimeImmutable $end,
        public string $city,
        public int $streetId,
        public string $streetName,
        /** @var string[] */
        public array $buildings,
        public string $comment,
    ) {
    }
}
