<?php

namespace App\Domain\ValueObject;

readonly class OutageData
{
    public function __construct(
        public \DateTimeImmutable $startDate,
        public \DateTimeImmutable $endDate,
        public string $comment,
    ) {
        if ($this->startDate > $this->endDate) {
            throw new \DomainException('Start date must be before or equal to end date');
        }
    }

    public function equals(OutageData $other): bool
    {
        return $this->startDate == $other->startDate
            && $this->endDate == $other->endDate
            && $this->comment === $other->comment;
    }
}
