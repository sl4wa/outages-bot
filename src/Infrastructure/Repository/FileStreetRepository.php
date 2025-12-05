<?php

declare(strict_types=1);

namespace App\Infrastructure\Repository;

use App\Application\Interface\Repository\StreetRepositoryInterface;
use App\Domain\Entity\Street;
use RuntimeException;
use Symfony\Component\DependencyInjection\ParameterBag\ParameterBagInterface;

final class FileStreetRepository implements StreetRepositoryInterface
{
    /** @var Street[] */
    private array $streets = [];

    public function __construct(ParameterBagInterface $params)
    {
        $projectDir = $params->get('kernel.project_dir');
        $streetsFile = $projectDir . '/data/streets.json';

        if (!file_exists($streetsFile)) {
            throw new RuntimeException('Streets file not found: ' . $streetsFile);
        }

        $json = file_get_contents($streetsFile);

        if ($json === false) {
            throw new RuntimeException('Failed to read streets file: ' . $streetsFile);
        }

        /** @var array<int, array{id: int, name: string}> $decoded */
        $decoded = json_decode($json, true) ?: [];

        foreach ($decoded as $st) {
            $this->streets[] = new Street($st['id'], $st['name']);
        }
    }

    public function getAllStreets(): array
    {
        return $this->streets;
    }
}
