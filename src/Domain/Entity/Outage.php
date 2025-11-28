<?php

declare(strict_types=1);

namespace App\Domain\Entity;

use App\Domain\ValueObject\OutageAddress;
use App\Domain\ValueObject\OutageDescription;
use App\Domain\ValueObject\OutagePeriod;
use App\Domain\ValueObject\UserAddress;

final readonly class Outage
{
    public function __construct(
        public int $id,
        public OutagePeriod $period,
        public OutageAddress $address,
        public OutageDescription $description,
    ) {
    }

    public function affectsUserAddress(UserAddress $userAddress): bool
    {
        return $this->address->coversUserAddress($userAddress);
    }
}
