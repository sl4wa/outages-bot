<?php
namespace App\Domain\Entity;

use App\Domain\ValueObject\Address;

readonly class Outage
{
    public function __construct(
        public \DateTimeImmutable $start,
        public \DateTimeImmutable $end,
        public Address $address,
        public string $comment
    ) {}

    public function isIdenticalPeriodAndComment(User $user): bool
    {
        return $user->startDate == $this->start
            && $user->endDate == $this->end
            && $user->comment == $this->comment;
    }
}
