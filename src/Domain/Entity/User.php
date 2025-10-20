<?php
namespace App\Domain\Entity;

use App\Domain\ValueObject\Address;

readonly class User
{
    public function __construct(
        public int $id,
        public Address $address,
        public ?\DateTimeImmutable $startDate,
        public ?\DateTimeImmutable $endDate,
        public string $comment,
    ) {}

    public function withUpdatedOutage(Outage $outage): self
    {
        return new self(
            $this->id,
            $this->address,
            $outage->start,
            $outage->end,
            $outage->comment
        );
    }
}
