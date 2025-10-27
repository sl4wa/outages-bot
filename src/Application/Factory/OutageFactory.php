<?php

namespace App\Application\Factory;

use App\Application\DTO\OutageDTO;
use App\Domain\Entity\Outage;
use App\Domain\ValueObject\Address;
use App\Domain\ValueObject\OutageData;

class OutageFactory
{
    public function createFromDTO(OutageDTO $dto): Outage
    {
        $outageData = new OutageData(
            $dto->start,
            $dto->end,
            $dto->comment
        );

        $address = new Address(
            $dto->streetId,
            $dto->streetName,
            $dto->buildings,
            $dto->city
        );

        return new Outage(
            $outageData,
            $address,
        );
    }
}

