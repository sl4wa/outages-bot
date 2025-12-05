<?php

declare(strict_types=1);

namespace App\Tests\Infrastructure\Repository;

use App\Infrastructure\Repository\FileStreetRepository;
use PHPUnit\Framework\TestCase;
use Symfony\Component\DependencyInjection\ParameterBag\ParameterBagInterface;

final class FileStreetRepositoryTest extends TestCase
{
    private FileStreetRepository $repository;

    protected function setUp(): void
    {
        $params = $this->createMock(ParameterBagInterface::class);
        $params->method('get')
            ->with('kernel.project_dir')
            ->willReturn(__DIR__ . '/../../fixtures');

        $this->repository = new FileStreetRepository($params);
    }

    public function testExactMatchReturnsSingleResult(): void
    {
        $result = $this->repository->filter('Сихівська');

        self::assertCount(1, $result);
        self::assertSame(1, $result[0]['id']);
        self::assertSame('Сихівська', $result[0]['name']);
    }

    public function testExactMatchCaseInsensitive(): void
    {
        $result = $this->repository->filter('сихівська');

        self::assertCount(1, $result);
        self::assertSame(1, $result[0]['id']);
        self::assertSame('Сихівська', $result[0]['name']);
    }

    public function testPartialMatchReturnsMultipleResults(): void
    {
        $result = $this->repository->filter('Сих');

        self::assertCount(2, $result);
    }

    public function testNoMatchReturnsEmptyArray(): void
    {
        $result = $this->repository->filter('Неіснуюча');

        self::assertCount(0, $result);
    }
}
