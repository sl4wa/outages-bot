<?php

declare(strict_types=1);

namespace App\Application\Interface\Repository;

interface StreetRepositoryInterface
{
    /**
     * Filter streets by query string (case-insensitive partial match).
     *
     * @return array<int, array{id: int, name: string}>
     */
    public function filter(string $query): array;

    /**
     * Find exact street by name (case-insensitive exact match).
     *
     * @return array{id: int, name: string}|null
     */
    public function findByName(string $name): ?array;
}
