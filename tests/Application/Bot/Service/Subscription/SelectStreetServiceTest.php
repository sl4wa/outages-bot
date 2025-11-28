<?php

declare(strict_types=1);

namespace App\Tests\Application\Bot\Service\Subscription;

use App\Application\Bot\Service\Subscription\SelectStreetService;
use App\Application\Interface\Repository\StreetRepositoryInterface;
use PHPUnit\Framework\TestCase;

final class SelectStreetServiceTest extends TestCase
{
    private SelectStreetService $service;
    private StreetRepositoryInterface $streetRepository;

    protected function setUp(): void
    {
        $this->streetRepository = $this->createMock(StreetRepositoryInterface::class);
        $this->service = new SelectStreetService($this->streetRepository);
    }

    public function testEmptyInputError(): void
    {
        $this->streetRepository
            ->expects($this->never())
            ->method('filter');

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
            ->method('filter');

        $result = $this->service->handle('   ');

        self::assertFalse($result->shouldContinue);
        self::assertSame('Введіть назву вулиці.', $result->message);
    }

    public function testNoStreetsFoundError(): void
    {
        $this->streetRepository
            ->expects($this->once())
            ->method('filter')
            ->with('Неіснуюча')
            ->willReturn([]);

        $this->streetRepository
            ->expects($this->never())
            ->method('findByName');

        $result = $this->service->handle('Неіснуюча');

        self::assertFalse($result->shouldContinue);
        self::assertSame('Вулицю не знайдено. Спробуйте ще раз.', $result->message);
        self::assertFalse($result->hasMultipleOptions());
        self::assertFalse($result->hasExactMatch());
    }

    public function testExactMatchFound(): void
    {
        $exactStreet = ['id' => 12783, 'name' => 'Шевченка Т.'];

        $this->streetRepository
            ->expects($this->once())
            ->method('filter')
            ->with('Шевченка Т.')
            ->willReturn([$exactStreet]);

        $this->streetRepository
            ->expects($this->once())
            ->method('findByName')
            ->with('Шевченка Т.')
            ->willReturn($exactStreet);

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
            ['id' => 12783, 'name' => 'Шевченка Т.'],
            ['id' => 12444, 'name' => 'Молдавська'],
            ['id' => 12445, 'name' => 'Стрийська'],
        ];

        $this->streetRepository
            ->expects($this->once())
            ->method('filter')
            ->with('вул')
            ->willReturn($streets);

        $this->streetRepository
            ->expects($this->once())
            ->method('findByName')
            ->with('вул')
            ->willReturn(null);

        $result = $this->service->handle('вул');

        self::assertFalse($result->shouldContinue);
        self::assertTrue($result->hasMultipleOptions());
        self::assertFalse($result->hasExactMatch());
        self::assertSame('Будь ласка, оберіть вулицю:', $result->message);
        self::assertCount(3, $result->streetOptions);
        self::assertSame($streets, $result->streetOptions);
    }

    public function testSinglePartialMatchWithoutExactMatch(): void
    {
        $streets = [
            ['id' => 12783, 'name' => 'Шевченка Т.'],
        ];

        $this->streetRepository
            ->expects($this->once())
            ->method('filter')
            ->with('Шевч')
            ->willReturn($streets);

        $this->streetRepository
            ->expects($this->once())
            ->method('findByName')
            ->with('Шевч')
            ->willReturn(null);

        $result = $this->service->handle('Шевч');

        self::assertFalse($result->shouldContinue);
        self::assertTrue($result->hasMultipleOptions());
        self::assertFalse($result->hasExactMatch());
        self::assertCount(1, $result->streetOptions);
    }

    public function testTrimsWhitespaceFromQuery(): void
    {
        $exactStreet = ['id' => 12783, 'name' => 'Шевченка Т.'];

        $this->streetRepository
            ->expects($this->once())
            ->method('filter')
            ->with('Шевченка Т.')
            ->willReturn([$exactStreet]);

        $this->streetRepository
            ->expects($this->once())
            ->method('findByName')
            ->with('Шевченка Т.')
            ->willReturn($exactStreet);

        $result = $this->service->handle('  Шевченка Т.  ');

        self::assertTrue($result->hasExactMatch());
        self::assertSame('Шевченка Т.', $result->selectedStreetName);
    }

    public function testHandlesCyrillicTextCorrectly(): void
    {
        $exactStreet = ['id' => 12444, 'name' => 'Молдавська'];

        $this->streetRepository
            ->expects($this->once())
            ->method('filter')
            ->with('Молдавська')
            ->willReturn([$exactStreet]);

        $this->streetRepository
            ->expects($this->once())
            ->method('findByName')
            ->with('Молдавська')
            ->willReturn($exactStreet);

        $result = $this->service->handle('Молдавська');

        self::assertTrue($result->hasExactMatch());
        self::assertSame('Молдавська', $result->selectedStreetName);
    }

    public function testExactMatchTakesPrecedenceOverMultipleResults(): void
    {
        $exactStreet = ['id' => 12783, 'name' => 'Київська'];
        $multipleStreets = [
            ['id' => 12783, 'name' => 'Київська'],
            ['id' => 99999, 'name' => 'Київська бічна'],
        ];

        $this->streetRepository
            ->expects($this->once())
            ->method('filter')
            ->with('Київська')
            ->willReturn($multipleStreets);

        $this->streetRepository
            ->expects($this->once())
            ->method('findByName')
            ->with('Київська')
            ->willReturn($exactStreet);

        $result = $this->service->handle('Київська');

        self::assertTrue($result->hasExactMatch());
        self::assertFalse($result->hasMultipleOptions());
        self::assertSame(12783, $result->selectedStreetId);
        self::assertStringContainsString('Ви обрали вулицю: Київська', $result->message);
    }
}
