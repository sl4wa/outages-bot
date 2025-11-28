<?php

declare(strict_types=1);

namespace App\Tests\Domain\ValueObject;

use App\Domain\ValueObject\OutagePeriod;
use PHPUnit\Framework\TestCase;

final class OutagePeriodTest extends TestCase
{
    public function testCreatesOutagePeriodWithValidDates(): void
    {
        $startDate = new \DateTimeImmutable('2024-11-28T06:47:00+00:00');
        $endDate = new \DateTimeImmutable('2024-11-28T10:00:00+00:00');

        $period = new OutagePeriod($startDate, $endDate);

        self::assertSame($startDate, $period->startDate);
        self::assertSame($endDate, $period->endDate);
        self::assertTrue($period->startDate <= $period->endDate);
    }

    public function testCreatesOutagePeriodWithEqualDates(): void
    {
        $date = new \DateTimeImmutable('2024-11-28T10:00:00+00:00');

        $period = new OutagePeriod($date, $date);

        self::assertSame($date, $period->startDate);
        self::assertSame($date, $period->endDate);
    }

    public function testThrowsExceptionWhenStartDateAfterEndDate(): void
    {
        $this->expectException(\DomainException::class);
        $this->expectExceptionMessage('Start date must be before or equal to end date');

        $startDate = new \DateTimeImmutable('2024-11-28T10:00:00+00:00');
        $endDate = new \DateTimeImmutable('2024-11-28T06:47:00+00:00');

        new OutagePeriod($startDate, $endDate);
    }

    public function testEqualsReturnsTrueForIdenticalPeriods(): void
    {
        $startDate = new \DateTimeImmutable('2024-11-28T06:47:00+00:00');
        $endDate = new \DateTimeImmutable('2024-11-28T10:00:00+00:00');

        $period1 = new OutagePeriod($startDate, $endDate);
        $period2 = new OutagePeriod($startDate, $endDate);

        self::assertTrue($period1->equals($period2));
    }

    public function testEqualsReturnsFalseForDifferentStartDates(): void
    {
        $startDate1 = new \DateTimeImmutable('2024-11-28T06:47:00+00:00');
        $startDate2 = new \DateTimeImmutable('2024-11-28T07:00:00+00:00');
        $endDate = new \DateTimeImmutable('2024-11-28T10:00:00+00:00');

        $period1 = new OutagePeriod($startDate1, $endDate);
        $period2 = new OutagePeriod($startDate2, $endDate);

        self::assertFalse($period1->equals($period2));
    }

    public function testEqualsReturnsFalseForDifferentEndDates(): void
    {
        $startDate = new \DateTimeImmutable('2024-11-28T06:47:00+00:00');
        $endDate1 = new \DateTimeImmutable('2024-11-28T10:00:00+00:00');
        $endDate2 = new \DateTimeImmutable('2024-11-28T11:00:00+00:00');

        $period1 = new OutagePeriod($startDate, $endDate1);
        $period2 = new OutagePeriod($startDate, $endDate2);

        self::assertFalse($period1->equals($period2));
    }

    public function testEqualsReturnsFalseForCompletelyDifferentPeriods(): void
    {
        $period1 = new OutagePeriod(
            new \DateTimeImmutable('2024-11-28T06:47:00+00:00'),
            new \DateTimeImmutable('2024-11-28T10:00:00+00:00')
        );
        $period2 = new OutagePeriod(
            new \DateTimeImmutable('2024-11-29T08:00:00+00:00'),
            new \DateTimeImmutable('2024-11-29T12:00:00+00:00')
        );

        self::assertFalse($period1->equals($period2));
    }
}
