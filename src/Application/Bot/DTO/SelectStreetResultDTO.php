<?php

namespace App\Application\Bot\DTO;

readonly class SelectStreetResultDTO
{
    public function __construct(
        public string $message,
        public ?array $streetOptions = null,
        public ?int $selectedStreetId = null,
        public ?string $selectedStreetName = null,
        public bool $shouldContinue = true
    ) {
    }

    public function hasMultipleOptions(): bool
    {
        return $this->streetOptions !== null && count($this->streetOptions) > 0;
    }

    public function hasExactMatch(): bool
    {
        return $this->selectedStreetId !== null && $this->selectedStreetName !== null;
    }
}
