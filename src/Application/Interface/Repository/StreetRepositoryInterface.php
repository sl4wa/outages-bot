<?php

declare(strict_types=1);

namespace App\Application\Interface\Repository;

interface StreetRepositoryInterface
{
    /**
     * @return array<int, array{id: int, name: string}>
     */
    public function filter(string $query): array;
}
