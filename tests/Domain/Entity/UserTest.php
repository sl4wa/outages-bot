<?php

declare(strict_types=1);

namespace App\Tests\Domain\Entity;

use App\Domain\Entity\Outage;
use App\Domain\Entity\User;
use App\Domain\ValueObject\OutageAddress;
use App\Domain\ValueObject\OutageDescription;
use App\Domain\ValueObject\OutageInfo;
use App\Domain\ValueObject\OutagePeriod;
use App\Domain\ValueObject\UserAddress;
use PHPUnit\Framework\TestCase;

final class UserTest extends TestCase
{
    private function createTestUserAddress(
        int $streetId = 12783,
        string $streetName = 'Шевченка Т.',
        string $building = '271'
    ): UserAddress {
        return new UserAddress(
            streetId: $streetId,
            streetName: $streetName,
            building: $building
        );
    }

    private function createTestOutage(
        int $id = 1,
        string $description = 'Застосування ГПВ'
    ): Outage {
        $period = new OutagePeriod(
            new \DateTimeImmutable('2024-11-28T06:47:00+00:00'),
            new \DateTimeImmutable('2024-11-28T10:00:00+00:00')
        );
        $address = new OutageAddress(
            streetId: 12783,
            streetName: 'Шевченка Т.',
            buildings: ['271', '273', '275'],
            city: 'Львів'
        );

        return new Outage($id, $period, $address, new OutageDescription($description));
    }

    public function testCreatesUserWithRequiredProperties(): void
    {
        $id = 123456;
        $address = $this->createTestUserAddress();

        $user = new User($id, $address, null);

        self::assertSame($id, $user->id);
        self::assertSame($address, $user->address);
        self::assertNull($user->outageInfo);
    }

    public function testCreatesUserWithOutageInfo(): void
    {
        $outageInfo = new OutageInfo(
            new OutagePeriod(
                new \DateTimeImmutable('2024-11-28T06:47:00+00:00'),
                new \DateTimeImmutable('2024-11-28T10:00:00+00:00')
            ),
            new OutageDescription('Застосування ГПВ')
        );

        $user = new User(123456, $this->createTestUserAddress(), $outageInfo);

        self::assertSame($outageInfo, $user->outageInfo);
    }

    public function testWithNotifiedOutageCreatesNewUserInstance(): void
    {
        $originalUser = new User(123456, $this->createTestUserAddress(), null);
        $outage = $this->createTestOutage();

        $newUser = $originalUser->withNotifiedOutage($outage);

        self::assertNotSame($originalUser, $newUser);
        self::assertNull($originalUser->outageInfo);
        self::assertNotNull($newUser->outageInfo);
    }

    public function testWithNotifiedOutageStoresCorrectOutageInfo(): void
    {
        $user = new User(123456, $this->createTestUserAddress(), null);
        $outage = $this->createTestOutage(1, 'Застосування ГПВ');

        $newUser = $user->withNotifiedOutage($outage);

        self::assertInstanceOf(OutageInfo::class, $newUser->outageInfo);
        self::assertTrue($newUser->outageInfo->period->equals($outage->period));
        self::assertTrue($newUser->outageInfo->description->equals($outage->description));
    }

    public function testWithNotifiedOutagePreservesUserIdAndAddress(): void
    {
        $id = 123456;
        $address = $this->createTestUserAddress(12783, 'Шевченка Т.', '271');
        $user = new User($id, $address, null);
        $outage = $this->createTestOutage();

        $newUser = $user->withNotifiedOutage($outage);

        self::assertSame($id, $newUser->id);
        self::assertSame($address, $newUser->address);
    }

    public function testIsAlreadyNotifiedAboutReturnsFalseWhenOutageInfoIsNull(): void
    {
        $user = new User(123456, $this->createTestUserAddress(), null);
        $outageInfo = new OutageInfo(
            new OutagePeriod(
                new \DateTimeImmutable('2024-11-28T06:47:00+00:00'),
                new \DateTimeImmutable('2024-11-28T10:00:00+00:00')
            ),
            new OutageDescription('Застосування ГПВ')
        );

        self::assertFalse($user->isAlreadyNotifiedAbout($outageInfo));
    }

    public function testIsAlreadyNotifiedAboutReturnsTrueForMatchingOutageInfo(): void
    {
        $outageInfo = new OutageInfo(
            new OutagePeriod(
                new \DateTimeImmutable('2024-11-28T06:47:00+00:00'),
                new \DateTimeImmutable('2024-11-28T10:00:00+00:00')
            ),
            new OutageDescription('Застосування ГПВ')
        );
        $user = new User(123456, $this->createTestUserAddress(), $outageInfo);

        self::assertTrue($user->isAlreadyNotifiedAbout($outageInfo));
    }

    public function testIsAlreadyNotifiedAboutReturnsFalseForDifferentOutageInfo(): void
    {
        $outageInfo1 = new OutageInfo(
            new OutagePeriod(
                new \DateTimeImmutable('2024-11-28T06:47:00+00:00'),
                new \DateTimeImmutable('2024-11-28T10:00:00+00:00')
            ),
            new OutageDescription('First description')
        );
        $outageInfo2 = new OutageInfo(
            new OutagePeriod(
                new \DateTimeImmutable('2024-11-28T06:47:00+00:00'),
                new \DateTimeImmutable('2024-11-28T10:00:00+00:00')
            ),
            new OutageDescription('Second description')
        );
        $user = new User(123456, $this->createTestUserAddress(), $outageInfo1);

        self::assertFalse($user->isAlreadyNotifiedAbout($outageInfo2));
    }
}
