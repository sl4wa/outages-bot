<?php

namespace App\Application\Factory;

use App\Application\DTO\OutageDTO;
use App\Domain\Entity\Outage;

class OutageFactory
{
    public function createFromDTO(OutageDTO $dto): Outage
    {
        return new Outage(
            $dto->start,
            $dto->end,
            $dto->city,
            $dto->streetId,
            $dto->streetName,
            $dto->buildingNames,
            $dto->comment,
        );
    }
}

