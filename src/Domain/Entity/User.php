<?php
namespace App\Domain\Entity;

class User
{
    public function __construct(
        public readonly int $id,
        public readonly int $streetId,
        public readonly string $streetName,
        public readonly string $building,
        public readonly ?\DateTimeImmutable $startDate,
        public readonly ?\DateTimeImmutable $endDate,
        public readonly string $comment,
    ) {}

    public function withUpdatedOutage(Outage $outage): self
    {
        return new self(
            $this->id,
            $this->streetId,
            $this->streetName,
            $this->building,
            $outage->start,
            $outage->end,
            $outage->comment
        );
    }

}
