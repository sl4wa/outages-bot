<?php

namespace App\Domain\ValueObject;

readonly class OutageDescription implements \JsonSerializable
{
    public function __construct(
        public string $value,
    ) {}

    public function equals(OutageDescription $other): bool
    {
        return $this->value === $other->value;
    }

    public function jsonSerialize(): string
    {
        return $this->value;
    }
}
