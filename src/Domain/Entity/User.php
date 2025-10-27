<?php
namespace App\Domain\Entity;

use App\Domain\ValueObject\Address;
use App\Domain\ValueObject\OutageData;

readonly class User
{
    public function __construct(
        public int $id,
        public Address $address,
        public ?OutageData $lastNotifiedOutage,
    ) {}

    public function withNotifiedOutage(OutageData $outageData): self
    {
        return new self(
            $this->id,
            $this->address,
            $outageData
        );
    }

    public function wasAlreadyNotifiedAbout(OutageData $outageData): bool
    {
        if ($this->lastNotifiedOutage === null) {
            return false;
        }

        return $this->lastNotifiedOutage->equals($outageData);
    }
}
