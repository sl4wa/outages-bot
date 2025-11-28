<?php

declare(strict_types=1);

namespace App\Tests\Domain\ValueObject;

use App\Domain\ValueObject\OutageDescription;
use App\Domain\ValueObject\OutageInfo;
use App\Domain\ValueObject\OutagePeriod;
use DateTimeImmutable;
use PHPUnit\Framework\TestCase;

final class OutageInfoTest extends TestCase
{
    private function createTestPeriod(
        string $start = '2024-11-28T06:47:00+00:00',
        string $end = '2024-11-28T10:00:00+00:00'
    ): OutagePeriod {
        return new OutagePeriod(
            new DateTimeImmutable($start),
            new DateTimeImmutable($end)
        );
    }

    private function createTestDescription(string $value = 'Застосування ГПВ'): OutageDescription
    {
        return new OutageDescription($value);
    }

    public function testCreatesOutageInfoWithPeriodAndDescription(): void
    {
        $period = $this->createTestPeriod();
        $description = $this->createTestDescription();

        $info = new OutageInfo($period, $description);

        self::assertSame($period, $info->period);
        self::assertSame($description, $info->description);
    }

    public function testEqualsReturnsTrueForIdenticalOutageInfo(): void
    {
        $period = $this->createTestPeriod();
        $description = $this->createTestDescription();

        $info1 = new OutageInfo($period, $description);
        $info2 = new OutageInfo($period, $description);

        self::assertTrue($info1->equals($info2));
    }

    public function testEqualsReturnsFalseForDifferentPeriods(): void
    {
        $period1 = $this->createTestPeriod('2024-11-28T06:47:00+00:00', '2024-11-28T10:00:00+00:00');
        $period2 = $this->createTestPeriod('2024-11-29T08:00:00+00:00', '2024-11-29T12:00:00+00:00');
        $description = $this->createTestDescription();

        $info1 = new OutageInfo($period1, $description);
        $info2 = new OutageInfo($period2, $description);

        self::assertFalse($info1->equals($info2));
    }

    public function testEqualsReturnsFalseForDifferentDescriptions(): void
    {
        $period = $this->createTestPeriod();
        $description1 = $this->createTestDescription('First description');
        $description2 = $this->createTestDescription('Second description');

        $info1 = new OutageInfo($period, $description1);
        $info2 = new OutageInfo($period, $description2);

        self::assertFalse($info1->equals($info2));
    }

    public function testEqualsReturnsFalseForCompletelyDifferentOutageInfo(): void
    {
        $period1 = $this->createTestPeriod('2024-11-28T06:47:00+00:00', '2024-11-28T10:00:00+00:00');
        $period2 = $this->createTestPeriod('2024-11-29T08:00:00+00:00', '2024-11-29T12:00:00+00:00');
        $description1 = $this->createTestDescription('First description');
        $description2 = $this->createTestDescription('Second description');

        $info1 = new OutageInfo($period1, $description1);
        $info2 = new OutageInfo($period2, $description2);

        self::assertFalse($info1->equals($info2));
    }
}
