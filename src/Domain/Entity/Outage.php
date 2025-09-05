<?php
namespace App\Domain\Entity;

class Outage
{
    public function __construct(
        public readonly \DateTimeImmutable $start,
        public readonly \DateTimeImmutable $end,
        public readonly string $city,
        public readonly int $streetId,
        public readonly string $streetName,
        public readonly array $buildingNames,
        public readonly string $comment
    ) {}

    public function matchesUser(User $user): bool
    {
        return $this->streetId === $user->streetId &&
            in_array($user->building, $this->buildingNames, true);
    }

    public function isIdenticalPeriodAndComment(User $user): bool
    {
        return $user->startDate == $this->start
            && $user->endDate == $this->end
            && $user->comment == $this->comment;
    }
}
