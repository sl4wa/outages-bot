<?php

namespace App\Domain\ValueObject;

readonly class UserAddress
{
    public function __construct(
        public int $streetId,
        public string $streetName,
        public string $building,
        public ?string $city = null,
    ) {
        if ($streetId <= 0) {
            throw new \InvalidArgumentException('Street ID must be positive');
        }

        if (trim($streetName) === '') {
            throw new \InvalidArgumentException('Street name cannot be empty');
        }

        if (trim($building) === '') {
            throw new \InvalidArgumentException('Building cannot be empty');
        }
    }
}
