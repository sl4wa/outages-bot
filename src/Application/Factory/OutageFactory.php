<?php

namespace App\Application\Factory;

use App\Application\DTO\OutageDTO;
use App\Domain\Entity\Outage;
use App\Domain\ValueObject\Address;

class OutageFactory
{
    public function createFromDTO(OutageDTO $dto): Outage
    {
        $address = new Address(
            $dto->streetId,
            $dto->streetName,
            $dto->buildingNames,
            $dto->city
        );

        return new Outage(
            $dto->start,
            $dto->end,
            $address,
            $dto->comment,
        );
    }
}

