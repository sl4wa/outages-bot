<?php

declare(strict_types=1);

namespace App\Tests\Domain\ValueObject;

use App\Domain\ValueObject\OutageAddress;
use App\Domain\ValueObject\UserAddress;
use InvalidArgumentException;
use PHPUnit\Framework\TestCase;

final class OutageAddressTest extends TestCase
{
    public function testCreatesOutageAddressWithValidData(): void
    {
        $address = new OutageAddress(
            streetId: 12783,
            streetName: 'Шевченка Т.',
            buildings: ['271', '273', '275']
        );

        self::assertSame(12783, $address->streetId);
        self::assertSame('Шевченка Т.', $address->streetName);
        self::assertSame(['271', '273', '275'], $address->buildings);
        self::assertNull($address->city);
    }

    public function testCreatesOutageAddressWithCity(): void
    {
        $address = new OutageAddress(
            streetId: 12783,
            streetName: 'Шевченка Т.',
            buildings: ['271'],
            city: 'Львів'
        );

        self::assertSame('Львів', $address->city);
    }

    public function testCoversUserAddressReturnsTrueForMatchingStreetAndBuilding(): void
    {
        $outageAddress = new OutageAddress(
            streetId: 12783,
            streetName: 'Шевченка Т.',
            buildings: ['271', '273', '275']
        );
        $userAddress = new UserAddress(
            streetId: 12783,
            streetName: 'Шевченка Т.',
            building: '271'
        );

        self::assertTrue($outageAddress->coversUserAddress($userAddress));
    }

    public function testCoversUserAddressReturnsFalseForDifferentStreet(): void
    {
        $outageAddress = new OutageAddress(
            streetId: 12783,
            streetName: 'Шевченка Т.',
            buildings: ['271']
        );
        $userAddress = new UserAddress(
            streetId: 99999,
            streetName: 'Інша вулиця',
            building: '271'
        );

        self::assertFalse($outageAddress->coversUserAddress($userAddress));
    }

    public function testCoversUserAddressReturnsFalseForNonMatchingBuilding(): void
    {
        $outageAddress = new OutageAddress(
            streetId: 12783,
            streetName: 'Шевченка Т.',
            buildings: ['271', '273']
        );
        $userAddress = new UserAddress(
            streetId: 12783,
            streetName: 'Шевченка Т.',
            building: '275'
        );

        self::assertFalse($outageAddress->coversUserAddress($userAddress));
    }

    public function testCoversUserAddressIsCaseSensitive(): void
    {
        $outageAddress = new OutageAddress(
            streetId: 12783,
            streetName: 'Шевченка Т.',
            buildings: ['271-А']  // Cyrillic А
        );
        $userAddress = new UserAddress(
            streetId: 12783,
            streetName: 'Шевченка Т.',
            building: '271-A'  // Latin A
        );

        self::assertFalse($outageAddress->coversUserAddress($userAddress));
    }

    /**
     * @dataProvider invalidStreetIdsProvider
     */
    public function testThrowsExceptionForNonPositiveStreetId(int $streetId): void
    {
        $this->expectException(InvalidArgumentException::class);
        $this->expectExceptionMessage('Street ID must be positive');

        new OutageAddress(
            streetId: $streetId,
            streetName: 'Шевченка Т.',
            buildings: ['271']
        );
    }

    public static function invalidStreetIdsProvider(): array
    {
        return [
            'zero' => [0],
            'negative' => [-1],
            'large negative' => [-100],
        ];
    }

    public function testThrowsExceptionForEmptyStreetName(): void
    {
        $this->expectException(InvalidArgumentException::class);
        $this->expectExceptionMessage('Street name cannot be empty');

        new OutageAddress(
            streetId: 12783,
            streetName: '',
            buildings: ['271']
        );
    }

    public function testThrowsExceptionForWhitespaceOnlyStreetName(): void
    {
        $this->expectException(InvalidArgumentException::class);
        $this->expectExceptionMessage('Street name cannot be empty');

        new OutageAddress(
            streetId: 12783,
            streetName: '   ',
            buildings: ['271']
        );
    }

    public function testThrowsExceptionForEmptyBuildingsArray(): void
    {
        $this->expectException(InvalidArgumentException::class);
        $this->expectExceptionMessage('Buildings must be non-empty strings');

        new OutageAddress(
            streetId: 12783,
            streetName: 'Шевченка Т.',
            buildings: []
        );
    }

    public function testThrowsExceptionForBuildingsArrayContainingEmptyString(): void
    {
        $this->expectException(InvalidArgumentException::class);
        $this->expectExceptionMessage('Buildings must be non-empty strings');

        new OutageAddress(
            streetId: 12783,
            streetName: 'Шевченка Т.',
            buildings: ['271', '', '273']
        );
    }

    public function testThrowsExceptionForBuildingsArrayContainingWhitespace(): void
    {
        $this->expectException(InvalidArgumentException::class);
        $this->expectExceptionMessage('Buildings must be non-empty strings');

        new OutageAddress(
            streetId: 12783,
            streetName: 'Шевченка Т.',
            buildings: ['271', '   ', '273']
        );
    }

    public function testThrowsExceptionForBuildingsArrayContainingNonString(): void
    {
        $this->expectException(InvalidArgumentException::class);
        $this->expectExceptionMessage('Buildings must be non-empty strings');

        new OutageAddress(
            streetId: 12783,
            streetName: 'Шевченка Т.',
            buildings: ['271', 123]
        );
    }
}
