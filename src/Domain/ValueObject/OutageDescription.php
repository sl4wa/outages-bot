<?php

namespace App\Domain\ValueObject;

readonly class OutageDescription
{
    public function __construct(
        public string $value,
    ) {}

    public function equals(OutageDescription $other): bool
    {
        return $this->value === $other->value;
    }

    public function isEmpty(): bool
    {
        return $this->value === '';
    }
}
