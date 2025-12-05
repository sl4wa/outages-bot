<?php

declare(strict_types=1);

namespace App\Application\Interface\Repository;

use App\Domain\Entity\Street;

interface StreetRepositoryInterface
{
    /**
     * @return Street[]
     */
    public function getAllStreets(): array;
}
