<?php

declare(strict_types=1);

namespace App\Infrastructure\Repository;

use App\Application\Interface\Repository\StreetRepositoryInterface;
use RuntimeException;
use Symfony\Component\DependencyInjection\ParameterBag\ParameterBagInterface;

final class FileStreetRepository implements StreetRepositoryInterface
{
    /** @var array<int, array{id: int, name: string, name_lower: string}> */
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
            $this->streets[] = [
                'id' => $st['id'],
                'name' => $st['name'],
                'name_lower' => mb_strtolower($st['name']),
            ];
        }
    }

    public function filter(string $query): array
    {
        $q = mb_strtolower(trim($query));

        $exactMatch = null;
        $results = [];

        foreach ($this->streets as $st) {
            if ($st['name_lower'] === $q) {
                $exactMatch = ['id' => $st['id'], 'name' => $st['name']];
                break;
            }

            if (str_contains($st['name_lower'], $q)) {
                $results[] = ['id' => $st['id'], 'name' => $st['name']];
            }
        }

        if ($exactMatch !== null) {
            return [$exactMatch];
        }

        return $results;
    }
}
