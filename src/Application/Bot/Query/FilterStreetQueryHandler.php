<?php

declare(strict_types=1);

namespace App\Application\Bot\Query;

use App\Application\Interface\Repository\StreetRepositoryInterface;
use App\Domain\Entity\Street;

final readonly class FilterStreetQueryHandler
{
    public function __construct(private StreetRepositoryInterface $streetRepository)
    {
    }

    /**
     * @return Street[]
     */
    public function handle(string $query): array
    {
        $q = mb_strtolower(trim($query));
        $streets = $this->streetRepository->getAllStreets();

        $results = [];

        foreach ($streets as $street) {
            if ($street->nameEquals($q)) {
                return [$street];
            }

            if ($street->nameContains($q)) {
                $results[] = $street;
            }
        }

        return $results;
    }
}
