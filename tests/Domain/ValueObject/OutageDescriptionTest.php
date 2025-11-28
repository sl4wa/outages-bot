<?php

declare(strict_types=1);

namespace App\Tests\Domain\ValueObject;

use App\Domain\ValueObject\OutageDescription;
use PHPUnit\Framework\TestCase;

final class OutageDescriptionTest extends TestCase
{
    public function testCreatesOutageDescriptionWithValue(): void
    {
        $description = new OutageDescription('Test description');

        self::assertSame('Test description', $description->value);
    }

    public function testCreatesOutageDescriptionWithEmptyString(): void
    {
        $description = new OutageDescription('');

        self::assertSame('', $description->value);
    }

    public function testCreatesOutageDescriptionWithCyrillicText(): void
    {
        $description = new OutageDescription('Застосування ГПВ');

        self::assertSame('Застосування ГПВ', $description->value);
    }

    public function testEqualsReturnsTrueForIdenticalDescriptions(): void
    {
        $description1 = new OutageDescription('Test description');
        $description2 = new OutageDescription('Test description');

        self::assertTrue($description1->equals($description2));
    }

    public function testEqualsReturnsFalseForDifferentDescriptions(): void
    {
        $description1 = new OutageDescription('First description');
        $description2 = new OutageDescription('Second description');

        self::assertFalse($description1->equals($description2));
    }

    public function testEqualsIsCaseSensitive(): void
    {
        $description1 = new OutageDescription('Test');
        $description2 = new OutageDescription('test');

        self::assertFalse($description1->equals($description2));
    }

    public function testJsonSerializeReturnsStringValue(): void
    {
        $description = new OutageDescription('Test description');

        $result = $description->jsonSerialize();

        self::assertIsString($result);
        self::assertSame('Test description', $result);
    }

    public function testJsonSerializeWithSpecialCharacters(): void
    {
        $value = "Description with \"quotes\" and \n newlines";
        $description = new OutageDescription($value);

        $result = $description->jsonSerialize();

        self::assertSame($value, $result);
    }
}
