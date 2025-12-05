<?php

declare(strict_types=1);

namespace App\Infrastructure\Repository;

use App\Application\Interface\Repository\StreetRepositoryInterface;
use RuntimeException;
use Symfony\Component\DependencyInjection\ParameterBag\ParameterBagInterface;

final class FileStreetRepository implements StreetRepositoryInterface
{
    private string $streetsFile;

    /** @var array<int, array{id: int, name: string}> */
    private array $streets = [];

    public function __construct(ParameterBagInterface $params)
    {
        $projectDir = $params->get('kernel.project_dir');
        $this->streetsFile = $projectDir . '/data/streets.json';

        if (!file_exists($this->streetsFile)) {
            throw new RuntimeException('Streets file not found: ' . $this->streetsFile);
        }

        $json = file_get_contents($this->streetsFile);

        if ($json === false) {
            throw new RuntimeException('Failed to read streets file: ' . $this->streetsFile);
        }

        /** @var array<int, array{id: int, name: string}> $decoded */
        $decoded = json_decode($json, true) ?: [];
        $this->streets = $decoded;
    }

    public function filter(string $query): array
    {
        $q = mb_strtolower(trim($query));

        return array_values(array_filter(
            $this->streets,
            fn ($st) => str_contains(mb_strtolower($st['name']), $q)
        ));
    }

    public function findByName(string $name): ?array
    {
        $norm = mb_strtolower(trim($name));

        foreach ($this->streets as $st) {
            if (mb_strtolower($st['name']) === $norm) {
                return $st;
            }
        }

        return null;
    }
}
