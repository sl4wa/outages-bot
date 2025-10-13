<?php
namespace App\Domain\Entity;

readonly class Outage
{
    public function __construct(
        public \DateTimeImmutable $start,
        public \DateTimeImmutable $end,
        public string $city,
        public int $streetId,
        public string $streetName,
        public array $buildingNames,
        public string $comment
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
