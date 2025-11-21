<?php

namespace App\Domain\ValueObject;

use App\Domain\Exception\InvalidBuildingFormatException;

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

        if (!preg_match('/^[0-9]+(-[A-Za-zА-Яа-яІіЇїЄєҐґ])?$/u', $building)) {
            throw new InvalidBuildingFormatException('Building format is invalid. Expected format: number or number-letter (e.g., 13 or 13-A)');
        }
    }
}
