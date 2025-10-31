<?php

namespace App\Domain\ValueObject;

readonly class OutageAddress
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

    public function coversUserAddress(UserAddress $userAddress): bool
    {
        return $this->streetId === $userAddress->streetId
            && in_array($userAddress->building, $this->buildings, true);
    }
}
