<?php
namespace App\Domain\Entity;

readonly class User
{
    public function __construct(
        public int $id,
        public int $streetId,
        public string $streetName,
        public string $building,
        public ?\DateTimeImmutable $startDate,
        public ?\DateTimeImmutable $endDate,
        public string $comment,
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
