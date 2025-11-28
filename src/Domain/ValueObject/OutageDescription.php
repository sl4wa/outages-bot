<?php

declare(strict_types=1);

namespace App\Domain\ValueObject;

use JsonSerializable;

final readonly class OutageDescription implements JsonSerializable
{
    public function __construct(
        public string $value,
    ) {
    }

    public function equals(self $other): bool
    {
        return $this->value === $other->value;
    }

    public function jsonSerialize(): string
    {
        return $this->value;
    }
}
