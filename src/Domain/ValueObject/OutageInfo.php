<?php

declare(strict_types=1);

namespace App\Domain\ValueObject;

final readonly class OutageInfo
{
    public function __construct(
        public OutagePeriod $period,
        public OutageDescription $description,
    ) {
    }

    public function equals(self $other): bool
    {
        return $this->period->equals($other->period)
            && $this->description->equals($other->description);
    }
}
