<?php

declare(strict_types=1);

namespace App\Tests\Domain\Entity;

use App\Domain\Entity\Outage;
use App\Domain\ValueObject\OutageAddress;
use App\Domain\ValueObject\OutageDescription;
use App\Domain\ValueObject\OutagePeriod;
use App\Domain\ValueObject\UserAddress;
use DateTimeImmutable;
use PHPUnit\Framework\TestCase;

final class OutageTest extends TestCase
{
    private function createTestPeriod(): OutagePeriod
    {
        return new OutagePeriod(
            new DateTimeImmutable('2024-11-28T06:47:00+00:00'),
            new DateTimeImmutable('2024-11-28T10:00:00+00:00')
        );
    }

    private function createTestAddress(
        int $streetId = 12783,
        string $streetName = 'Шевченка Т.',
        array $buildings = ['271', '273', '275']
    ): OutageAddress {
        return new OutageAddress(
            streetId: $streetId,
            streetName: $streetName,
            buildings: $buildings,
            city: 'Львів'
        );
    }

    private function createTestDescription(string $value = 'Застосування ГПВ'): OutageDescription
    {
        return new OutageDescription($value);
    }

    public function testCreatesOutageWithAllProperties(): void
    {
        $id = 170149994;
        $period = $this->createTestPeriod();
        $address = $this->createTestAddress();
        $description = $this->createTestDescription();

        $outage = new Outage($id, $period, $address, $description);

        self::assertSame($id, $outage->id);
        self::assertSame($period, $outage->period);
        self::assertSame($address, $outage->address);
        self::assertSame($description, $outage->description);
    }

    public function testAffectsUserAddressReturnsTrueForMatchingAddress(): void
    {
        $outage = new Outage(
            id: 1,
            period: $this->createTestPeriod(),
            address: $this->createTestAddress(12783, 'Шевченка Т.', ['271', '273', '275']),
            description: $this->createTestDescription()
        );

        $userAddress = new UserAddress(
            streetId: 12783,
            streetName: 'Шевченка Т.',
            building: '271'
        );

        self::assertTrue($outage->affectsUserAddress($userAddress));
    }

    public function testAffectsUserAddressReturnsFalseForDifferentStreet(): void
    {
        $outage = new Outage(
            id: 1,
            period: $this->createTestPeriod(),
            address: $this->createTestAddress(12783, 'Шевченка Т.', ['271']),
            description: $this->createTestDescription()
        );

        $userAddress = new UserAddress(
            streetId: 99999,
            streetName: 'Інша вулиця',
            building: '271'
        );

        self::assertFalse($outage->affectsUserAddress($userAddress));
    }

    public function testAffectsUserAddressReturnsFalseForNonCoveredBuilding(): void
    {
        $outage = new Outage(
            id: 1,
            period: $this->createTestPeriod(),
            address: $this->createTestAddress(12783, 'Шевченка Т.', ['271', '273']),
            description: $this->createTestDescription()
        );

        $userAddress = new UserAddress(
            streetId: 12783,
            streetName: 'Шевченка Т.',
            building: '275'
        );

        self::assertFalse($outage->affectsUserAddress($userAddress));
    }

    public function testAffectsUserAddressReturnsTrueForBuildingWithLetter(): void
    {
        $outage = new Outage(
            id: 1,
            period: $this->createTestPeriod(),
            address: $this->createTestAddress(12783, 'Шевченка Т.', ['271-А', '273', '275']),
            description: $this->createTestDescription()
        );

        $userAddress = new UserAddress(
            streetId: 12783,
            streetName: 'Шевченка Т.',
            building: '271-А'
        );

        self::assertTrue($outage->affectsUserAddress($userAddress));
    }
}
