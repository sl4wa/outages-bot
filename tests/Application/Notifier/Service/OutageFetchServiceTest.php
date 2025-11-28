<?php

declare(strict_types=1);

namespace App\Tests\Application\Notifier\Service;

use App\Application\Notifier\DTO\OutageDTO;
use App\Application\Notifier\Interface\Provider\OutageProviderInterface;
use App\Application\Notifier\Service\OutageFetchService;
use App\Domain\Entity\Outage;
use PHPUnit\Framework\TestCase;

final class OutageFetchServiceTest extends TestCase
{
    public function testHandleReturnsEmptyArrayWhenProviderReturnsEmpty(): void
    {
        $mockProvider = $this->createMock(OutageProviderInterface::class);
        $mockProvider->expects($this->once())
            ->method('fetchOutages')
            ->willReturn([]);

        $service = new OutageFetchService($mockProvider);
        $result = $service->handle();

        self::assertIsArray($result);
        self::assertEmpty($result);
    }

    public function testHandleMapsOutageDTOToOutageEntity(): void
    {
        $dto = $this->createTestOutageDTO(
            id: 170149994,
            streetName: 'Молдавська',
            streetId: 12444,
            buildings: ['1', '10', '13-А'],
            comment: 'Застосування ГПВ'
        );

        $mockProvider = $this->createMock(OutageProviderInterface::class);
        $mockProvider->expects($this->once())
            ->method('fetchOutages')
            ->willReturn([$dto]);

        $service = new OutageFetchService($mockProvider);
        $result = $service->handle();

        self::assertCount(1, $result);
        self::assertInstanceOf(Outage::class, $result[0]);
    }

    public function testHandleMapsMultipleOutageDTOs(): void
    {
        $dto1 = $this->createTestOutageDTO(id: 1);
        $dto2 = $this->createTestOutageDTO(id: 2);
        $dto3 = $this->createTestOutageDTO(id: 3);

        $mockProvider = $this->createMock(OutageProviderInterface::class);
        $mockProvider->expects($this->once())
            ->method('fetchOutages')
            ->willReturn([$dto1, $dto2, $dto3]);

        $service = new OutageFetchService($mockProvider);
        $result = $service->handle();

        self::assertCount(3, $result);
        self::assertSame(1, $result[0]->id);
        self::assertSame(2, $result[1]->id);
        self::assertSame(3, $result[2]->id);
    }

    public function testHandleCreatesOutagePeriodFromDTODates(): void
    {
        $dto = $this->createTestOutageDTO(
            startDate: '2024-11-28T06:47:00+00:00',
            endDate: '2024-11-28T10:00:00+00:00'
        );

        $mockProvider = $this->createMock(OutageProviderInterface::class);
        $mockProvider->expects($this->once())
            ->method('fetchOutages')
            ->willReturn([$dto]);

        $service = new OutageFetchService($mockProvider);
        $result = $service->handle();

        $outage = $result[0];
        self::assertEquals(new \DateTimeImmutable('2024-11-28T06:47:00+00:00'), $outage->period->startDate);
        self::assertEquals(new \DateTimeImmutable('2024-11-28T10:00:00+00:00'), $outage->period->endDate);
    }

    public function testHandleCreatesOutageAddressFromDTOData(): void
    {
        $dto = $this->createTestOutageDTO(
            streetName: 'Молдавська',
            streetId: 12444,
            buildings: ['1', '13-А', '76-Г'],
            city: 'Львів'
        );

        $mockProvider = $this->createMock(OutageProviderInterface::class);
        $mockProvider->expects($this->once())
            ->method('fetchOutages')
            ->willReturn([$dto]);

        $service = new OutageFetchService($mockProvider);
        $result = $service->handle();

        $outage = $result[0];
        self::assertSame(12444, $outage->address->streetId);
        self::assertSame('Молдавська', $outage->address->streetName);
        self::assertSame(['1', '13-А', '76-Г'], $outage->address->buildings);
        self::assertSame('Львів', $outage->address->city);
    }

    public function testHandleCreatesOutageDescriptionFromDTOComment(): void
    {
        $dto = $this->createTestOutageDTO(
            comment: 'Застосування ГПВ'
        );

        $mockProvider = $this->createMock(OutageProviderInterface::class);
        $mockProvider->expects($this->once())
            ->method('fetchOutages')
            ->willReturn([$dto]);

        $service = new OutageFetchService($mockProvider);
        $result = $service->handle();

        $outage = $result[0];
        self::assertSame('Застосування ГПВ', $outage->description->value);
    }

    public function testHandlePreservesOutageIdFromDTO(): void
    {
        $dto = $this->createTestOutageDTO(id: 170149994);

        $mockProvider = $this->createMock(OutageProviderInterface::class);
        $mockProvider->expects($this->once())
            ->method('fetchOutages')
            ->willReturn([$dto]);

        $service = new OutageFetchService($mockProvider);
        $result = $service->handle();

        $outage = $result[0];
        self::assertSame(170149994, $outage->id);
    }

    public function testHandleWorksWithRealDataStructure(): void
    {
        $dto = $this->createTestOutageDTO(
            id: 170149994,
            streetName: 'Молдавська',
            streetId: 12444,
            buildings: [
                '1', '10', '11', '13', '13-А', '13-Б', '15', '15-А', '15-Б', '16',
                '17-А', '17-В', '19', '2', '21', '22', '23', '23-А', '25', '26',
                '27', '28', '29', '3', '30', '31', '31-А', '33', '33-А', '34',
                '35', '36', '37', '39', '39-Б', '39-Г', '4', '40', '43', '45',
                '47', '49', '5', '51', '53', '55', '57', '59', '6', '6-А',
                '61', '7', '76', '76-А', '76-Б', '76-В', '76-Г', '76-Д', '76-Е',
                '8', '86', '9'
            ],
            comment: 'Застосування ГПВ',
            city: 'Львів',
            startDate: '2025-11-26T08:00:00+00:00',
            endDate: '2025-11-26T10:30:00+00:00'
        );

        $mockProvider = $this->createMock(OutageProviderInterface::class);
        $mockProvider->expects($this->once())
            ->method('fetchOutages')
            ->willReturn([$dto]);

        $service = new OutageFetchService($mockProvider);
        $result = $service->handle();

        self::assertCount(1, $result);
        $outage = $result[0];

        self::assertSame(170149994, $outage->id);
        self::assertSame('Молдавська', $outage->address->streetName);
        self::assertSame(62, count($outage->address->buildings));
        self::assertContains('13-А', $outage->address->buildings);
        self::assertContains('76-Г', $outage->address->buildings);
        self::assertSame('Застосування ГПВ', $outage->description->value);
    }

    private function createTestOutageDTO(
        int $id = 1,
        string $streetName = 'Шевченка Т.',
        int $streetId = 12783,
        array $buildings = ['271', '273', '275'],
        string $comment = 'Застосування ГПВ',
        string $city = 'Львів',
        string $startDate = '2024-11-28T06:47:00+00:00',
        string $endDate = '2024-11-28T10:00:00+00:00'
    ): OutageDTO {
        return new OutageDTO(
            id: $id,
            start: new \DateTimeImmutable($startDate),
            end: new \DateTimeImmutable($endDate),
            city: $city,
            streetId: $streetId,
            streetName: $streetName,
            buildings: $buildings,
            comment: $comment
        );
    }
}
