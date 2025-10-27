<?php
namespace App\Domain\Entity;

use App\Domain\ValueObject\Address;
use App\Domain\ValueObject\OutageData;

readonly class Outage
{
    public function __construct(
        public OutageData $data,
        public Address $address,
    ) {}
}
