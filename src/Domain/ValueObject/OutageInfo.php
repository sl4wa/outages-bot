<?php

namespace App\Domain\ValueObject;

readonly class OutageInfo
{
    public function __construct(
        public OutagePeriod $period,
        public OutageDescription $description,
    ) {}

    public function equals(OutageInfo $other): bool
    {
        return $this->period->equals($other->period)
            && $this->description->equals($other->description);
    }
}
