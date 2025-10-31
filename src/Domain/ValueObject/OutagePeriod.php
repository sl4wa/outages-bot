<?php

namespace App\Domain\ValueObject;

readonly class OutagePeriod
{
    public function __construct(
        public \DateTimeImmutable $startDate,
        public \DateTimeImmutable $endDate,
    ) {
        if ($this->startDate > $this->endDate) {
            throw new \DomainException('Start date must be before or equal to end date');
        }
    }

    public function equals(OutagePeriod $other): bool
    {
        return $this->startDate == $other->startDate
            && $this->endDate == $other->endDate;
    }
}

