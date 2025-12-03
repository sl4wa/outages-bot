<?php

declare(strict_types=1);

namespace App\Tests\Domain\ValueObject;

use App\Domain\Exception\InvalidBuildingFormatException;
use App\Domain\ValueObject\UserAddress;
use InvalidArgumentException;
use PHPUnit\Framework\TestCase;

final class UserAddressTest extends TestCase
{
    public function testCreatesUserAddressWithValidData(): void
    {
        $address = new UserAddress(
            streetId: 123,
            streetName: 'Шевченка',
            building: '13'
        );

        self::assertSame(123, $address->streetId);
        self::assertSame('Шевченка', $address->streetName);
        self::assertSame('13', $address->building);
        self::assertNull($address->city);
    }

    public function testCreatesUserAddressWithCity(): void
    {
        $address = new UserAddress(
            streetId: 123,
            streetName: 'Шевченка',
            building: '13',
            city: 'Львів'
        );

        self::assertSame('Львів', $address->city);
    }

    /**
     * @dataProvider validBuildingFormatsProvider
     */
    public function testAcceptsValidBuildingFormats(string $building): void
    {
        $address = new UserAddress(
            streetId: 123,
            streetName: 'Шевченка',
            building: $building
        );

        self::assertSame($building, $address->building);
    }

    public static function validBuildingFormatsProvider(): array
    {
        return [
            'simple number' => ['13'],
            'large number' => ['196'],
            'three digit number' => ['271'],
            'number with latin letter' => ['13-A'],
            'number with cyrillic А' => ['196-А'],
            'number with cyrillic Б' => ['271-Б'],
            'number with cyrillic В' => ['350-В'],
            'number with cyrillic І' => ['25-І'],
            'number with cyrillic Ї' => ['30-Ї'],
            'number with cyrillic Є' => ['40-Є'],
            'number with cyrillic Ґ' => ['50-Ґ'],
        ];
    }

    /**
     * @dataProvider invalidBuildingFormatsProvider
     */
    public function testRejectsInvalidBuildingFormats(string $building): void
    {
        $this->expectException(InvalidBuildingFormatException::class);
        $this->expectExceptionMessage('Building format is invalid. Expected format: number or number-letter (e.g., 13 or 13-A)');

        new UserAddress(
            streetId: 123,
            streetName: 'Шевченка',
            building: $building
        );
    }

    public static function invalidBuildingFormatsProvider(): array
    {
        return [
            'with parentheses' => ['291(0083)'],
            'with comma' => ['350,А'],
            'multiple letters' => ['13-AB'],
            'slash separator' => ['13/A'],
            'only letters' => ['abc'],
            'letter before number' => ['A-13'],
            'number after hyphen' => ['13-1'],
            'multiple hyphens' => ['13-A-B'],
            'space separator' => ['13 A'],
            'special characters' => ['13@A'],
            'dot separator' => ['13.A'],
            'no hyphen before letter' => ['13A'],
            'hyphen without letter' => ['13-'],
            'starts with hyphen' => ['-13'],
            'range with hyphen' => ['59-61'],
            'range with slash' => ['59/61'],
            'long range with hyphen' => ['125-131'],
            'spaced hyphen with cyrillic letter' => ['64 - А'],
            'fraction format' => ['180/4'],
            'lowercase latin letter' => ['13-a'],
            'lowercase cyrillic letter' => ['13-а'],
        ];
    }

    public function testThrowsExceptionForEmptyBuilding(): void
    {
        $this->expectException(InvalidArgumentException::class);
        $this->expectExceptionMessage('Building cannot be empty');

        new UserAddress(
            streetId: 123,
            streetName: 'Шевченка',
            building: ''
        );
    }

    public function testThrowsExceptionForWhitespaceOnlyBuilding(): void
    {
        $this->expectException(InvalidArgumentException::class);
        $this->expectExceptionMessage('Building cannot be empty');

        new UserAddress(
            streetId: 123,
            streetName: 'Шевченка',
            building: '   '
        );
    }

    public function testThrowsExceptionForInvalidStreetId(): void
    {
        $this->expectException(InvalidArgumentException::class);
        $this->expectExceptionMessage('Street ID must be positive');

        new UserAddress(
            streetId: 0,
            streetName: 'Шевченка',
            building: '13'
        );
    }

    public function testThrowsExceptionForNegativeStreetId(): void
    {
        $this->expectException(InvalidArgumentException::class);
        $this->expectExceptionMessage('Street ID must be positive');

        new UserAddress(
            streetId: -1,
            streetName: 'Шевченка',
            building: '13'
        );
    }

    public function testThrowsExceptionForEmptyStreetName(): void
    {
        $this->expectException(InvalidArgumentException::class);
        $this->expectExceptionMessage('Street name cannot be empty');

        new UserAddress(
            streetId: 123,
            streetName: '',
            building: '13'
        );
    }

    public function testThrowsExceptionForWhitespaceOnlyStreetName(): void
    {
        $this->expectException(InvalidArgumentException::class);
        $this->expectExceptionMessage('Street name cannot be empty');

        new UserAddress(
            streetId: 123,
            streetName: '   ',
            building: '13'
        );
    }
}
