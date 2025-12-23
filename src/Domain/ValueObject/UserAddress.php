<?php

declare(strict_types=1);

namespace App\Domain\ValueObject;

use InvalidArgumentException;

final readonly class UserAddress
{
    public function __construct(
        public int $streetId,
        public string $streetName,
        public string $building,
        public ?string $city = null,
    ) {
        if ($streetId <= 0) {
            throw new InvalidArgumentException('Невірний ідентифікатор вулиці');
        }

        if (trim($streetName) === '') {
            throw new InvalidArgumentException('Назва вулиці не може бути порожньою');
        }

        if (trim($building) === '') {
            throw new InvalidArgumentException('Невірний формат номера будинку');
        }

        if (!preg_match('/^[0-9]+(-[A-ZА-ЯІЇЄҐ])?$/u', $building)) {
            throw new InvalidArgumentException('Невірний формат номера будинку. Приклад: 13 або 13-А');
        }
    }
}
