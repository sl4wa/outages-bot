<?php

declare(strict_types=1);

namespace App\Domain\ValueObject;

use DateTimeImmutable;
use DomainException;

final readonly class OutagePeriod
{
    public function __construct(
        public DateTimeImmutable $startDate,
        public DateTimeImmutable $endDate,
    ) {
        if ($this->startDate > $this->endDate) {
            throw new DomainException('Start date must be before or equal to end date');
        }
    }

    public function equals(self $other): bool
    {
        return $this->startDate->getTimestamp() === $other->startDate->getTimestamp()
            && $this->endDate->getTimestamp() === $other->endDate->getTimestamp();
    }
}
