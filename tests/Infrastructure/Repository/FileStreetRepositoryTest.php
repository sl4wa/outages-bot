<?php

declare(strict_types=1);

namespace App\Tests\Infrastructure\Repository;

use App\Domain\Entity\Street;
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

    public function testGetAllStreetsReturnsAllStreets(): void
    {
        $result = $this->repository->getAllStreets();

        self::assertCount(2, $result);
        self::assertContainsOnlyInstancesOf(Street::class, $result);
    }

    public function testGetAllStreetsReturnsStreetEntities(): void
    {
        $result = $this->repository->getAllStreets();

        self::assertSame(1, $result[0]->id);
        self::assertSame('Сихівська', $result[0]->name);
    }
}
