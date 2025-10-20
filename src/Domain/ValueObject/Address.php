<?php

namespace App\Domain\ValueObject;

readonly class Address
{
    /** @param string[] $buildings */
    public function __construct(
        public int $streetId,
        public string $streetName,
        public array $buildings,
        public ?string $city = null,
    ) {
        if ($streetId <= 0) {
            throw new \InvalidArgumentException('Street ID must be positive');
        }

        if (trim($streetName) === '') {
            throw new \InvalidArgumentException('Street name cannot be empty');
        }

        if (empty($buildings) || array_filter($buildings, fn($b) => !is_string($b) || trim($b) === '')) {
            throw new \InvalidArgumentException('Buildings must be non-empty strings');
        }
    }

    public function covers(Address $other): bool
    {
        return $this->streetId === $other->streetId
            && count($other->buildings) === 1
            && in_array($other->buildings[0], $this->buildings, true);
    }

    public function getSingleBuilding(): string
    {
        if (count($this->buildings) !== 1) {
            throw new \LogicException('Address must have exactly one building');
        }

        return $this->buildings[0];
    }
}
