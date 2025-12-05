<?php

declare(strict_types=1);

namespace App\Domain\Entity;

final readonly class Street
{
    public function __construct(
        public int $id,
        public string $name,
    ) {
    }

    public function nameContains(string $query): bool
    {
        return str_contains(mb_strtolower($this->name), $query);
    }

    public function nameEquals(string $query): bool
    {
        return mb_strtolower($this->name) === $query;
    }
}
