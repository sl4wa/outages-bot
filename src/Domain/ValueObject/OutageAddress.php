<?php

declare(strict_types=1);

namespace App\Domain\ValueObject;

use InvalidArgumentException;

final readonly class OutageAddress
{
    /** @param string[] $buildings */
    public function __construct(
        public int $streetId,
        public string $streetName,
        public array $buildings,
        public ?string $city = null,
    ) {
        if ($streetId <= 0) {
            throw new InvalidArgumentException('Street ID must be positive');
        }

        if (trim($streetName) === '') {
            throw new InvalidArgumentException('Street name cannot be empty');
        }

        if (empty($buildings)) {
            throw new InvalidArgumentException('Buildings must be non-empty strings');
        }

        foreach ($buildings as $building) {
            // @phpstan-ignore function.alreadyNarrowedType (runtime validation for mixed input)
            if (!is_string($building) || trim($building) === '') {
                throw new InvalidArgumentException('Buildings must be non-empty strings');
            }
        }
    }

    public function coversUserAddress(UserAddress $userAddress): bool
    {
        return $this->streetId === $userAddress->streetId
            && in_array($userAddress->building, $this->buildings, true);
    }
}
