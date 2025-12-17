<?php

declare(strict_types=1);

namespace App\Tests\Application\Bot\Service;

use App\Application\Bot\Query\FilterStreetQueryHandler;
use App\Application\Bot\Service\SearchStreetService;
use App\Application\Interface\Repository\StreetRepositoryInterface;
use App\Domain\Entity\Street;
use PHPUnit\Framework\TestCase;

final class SearchStreetServiceTest extends TestCase
{
    private SearchStreetService $service;

    private StreetRepositoryInterface $streetRepository;

    protected function setUp(): void
    {
        $this->streetRepository = $this->createMock(StreetRepositoryInterface::class);
        $queryHandler = new FilterStreetQueryHandler($this->streetRepository);
        $this->service = new SearchStreetService($queryHandler);
    }

    public function testEmptyInputError(): void
    {
        $this->streetRepository
            ->expects($this->never())
            ->method('getAllStreets');

        $result = $this->service->handle('');

        self::assertFalse($result->shouldContinue);
        self::assertSame('Введіть назву вулиці.', $result->message);
        self::assertFalse($result->hasMultipleOptions());
        self::assertFalse($result->hasExactMatch());
    }

    public function testWhitespaceOnlyInputError(): void
    {
        $this->streetRepository
            ->expects($this->never())
            ->method('getAllStreets');

        $result = $this->service->handle('   ');

        self::assertFalse($result->shouldContinue);
        self::assertSame('Введіть назву вулиці.', $result->message);
    }

    public function testNoStreetsFoundError(): void
    {
        $this->streetRepository
            ->expects($this->once())
            ->method('getAllStreets')
            ->willReturn([]);

        $result = $this->service->handle('Неіснуюча');

        self::assertFalse($result->shouldContinue);
        self::assertSame('Вулицю не знайдено. Спробуйте ще раз.', $result->message);
        self::assertFalse($result->hasMultipleOptions());
        self::assertFalse($result->hasExactMatch());
    }

    public function testSingleMatchAutoSelected(): void
    {
        $street = new Street(id: 12783, name: 'Шевченка Т.');

        $this->streetRepository
            ->expects($this->once())
            ->method('getAllStreets')
            ->willReturn([$street]);

        $result = $this->service->handle('Шевченка Т.');

        self::assertTrue($result->shouldContinue);
        self::assertTrue($result->hasExactMatch());
        self::assertFalse($result->hasMultipleOptions());
        self::assertSame(12783, $result->selectedStreetId);
        self::assertSame('Шевченка Т.', $result->selectedStreetName);
        self::assertStringContainsString('Ви обрали вулицю: Шевченка Т.', $result->message);
        self::assertStringContainsString('введіть номер будинку', $result->message);
    }

    public function testMultipleMatchesFound(): void
    {
        $streets = [
            new Street(id: 12783, name: 'вул. Шевченка'),
            new Street(id: 12444, name: 'вул. Молдавська'),
            new Street(id: 12445, name: 'вул. Стрийська'),
        ];

        $this->streetRepository
            ->expects($this->once())
            ->method('getAllStreets')
            ->willReturn($streets);

        $result = $this->service->handle('вул');

        self::assertFalse($result->shouldContinue);
        self::assertTrue($result->hasMultipleOptions());
        self::assertFalse($result->hasExactMatch());
        self::assertSame('Будь ласка, оберіть вулицю:', $result->message);
        self::assertCount(3, $result->streetOptions);
    }

    public function testSinglePartialMatchAutoSelected(): void
    {
        $street = new Street(id: 12783, name: 'Шевченка Т.');

        $this->streetRepository
            ->expects($this->once())
            ->method('getAllStreets')
            ->willReturn([$street]);

        $result = $this->service->handle('Шевч');

        self::assertTrue($result->shouldContinue);
        self::assertTrue($result->hasExactMatch());
        self::assertSame(12783, $result->selectedStreetId);
        self::assertSame('Шевченка Т.', $result->selectedStreetName);
    }

    public function testTrimsWhitespaceFromQuery(): void
    {
        $street = new Street(id: 12783, name: 'Шевченка Т.');

        $this->streetRepository
            ->expects($this->once())
            ->method('getAllStreets')
            ->willReturn([$street]);

        $result = $this->service->handle('  Шевченка Т.  ');

        self::assertTrue($result->hasExactMatch());
        self::assertSame('Шевченка Т.', $result->selectedStreetName);
    }

    public function testHandlesCyrillicTextCorrectly(): void
    {
        $street = new Street(id: 12444, name: 'Молдавська');

        $this->streetRepository
            ->expects($this->once())
            ->method('getAllStreets')
            ->willReturn([$street]);

        $result = $this->service->handle('Молдавська');

        self::assertTrue($result->hasExactMatch());
        self::assertSame('Молдавська', $result->selectedStreetName);
    }

    public function testMultipleResultsShowOptions(): void
    {
        $streets = [
            new Street(id: 12783, name: 'Київська основна'),
            new Street(id: 99999, name: 'Київська бічна'),
        ];

        $this->streetRepository
            ->expects($this->once())
            ->method('getAllStreets')
            ->willReturn($streets);

        $result = $this->service->handle('Київська');

        self::assertFalse($result->hasExactMatch());
        self::assertTrue($result->hasMultipleOptions());
        self::assertCount(2, $result->streetOptions);
        self::assertSame('Будь ласка, оберіть вулицю:', $result->message);
    }
}
