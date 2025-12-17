<?php

declare(strict_types=1);

namespace App\Application\Bot\DTO;

use App\Domain\Entity\Street;

final readonly class SearchStreetResultDTO
{
    /**
     * @param Street[]|null $streetOptions
     */
    public function __construct(
        public string $message,
        public ?array $streetOptions = null,
        public ?int $selectedStreetId = null,
        public ?string $selectedStreetName = null,
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
