<?php

declare(strict_types=1);

namespace App\Domain\Repository;

use App\Domain\Entity\Street;

interface StreetRepositoryInterface
{
    /**
     * @return Street[]
     */
    public function getAllStreets(): array;
}
