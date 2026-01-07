<?php

declare(strict_types=1);

namespace App\Infrastructure\Provider;

use App\Application\Interface\OutageProviderInterface;
use App\Application\Notifier\DTO\OutageDTO;
use DateTimeImmutable;
use Symfony\Contracts\HttpClient\HttpClientInterface;

final class ApiOutageProvider implements OutageProviderInterface
{
    private const API_URL = 'https://power-api.loe.lviv.ua/api/pw_accidents?pagination=false&otg.id=28&city.id=693';

    public function __construct(private readonly HttpClientInterface $httpClient)
    {
    }

    /**
     * @return OutageDTO[]
     */
    public function fetchOutages(): array
    {
        $response = $this->httpClient->request('GET', self::API_URL);

        if ($response->getStatusCode() !== 200) {
            return [];
        }
        $data = $response->toArray();
        /** @var array<int, array<string, mixed>> $items */
        $items = $data['hydra:member'] ?? [];

        $result = [];

        foreach ($items as $row) {
            $id = $this->getInt($row, 'id');

            $comment = $this->getString($row, 'koment');
            $comment = preg_replace('/[\r\n]+/', ' ', $comment) ?? '';
            $comment = trim($comment);

            $buildingsRaw = $row['buildingNames'] ?? '';
            /** @var string[] $buildings */
            $buildings = is_array($buildingsRaw)
                ? array_map(fn (mixed $b): string => trim($this->castToString($b)), $buildingsRaw)
                : array_values(array_filter(array_map('trim', explode(',', $this->castToString($buildingsRaw)))));

            /** @var array<string, mixed> $city */
            $city = is_array($row['city'] ?? null) ? $row['city'] : [];
            /** @var array<string, mixed> $street */
            $street = is_array($row['street'] ?? null) ? $row['street'] : [];

            $result[$id] = new OutageDTO(
                $id,
                new DateTimeImmutable($this->getString($row, 'dateEvent') ?: 'now'),
                new DateTimeImmutable($this->getString($row, 'datePlanIn') ?: 'now'),
                $this->getString($city, 'name'),
                $this->getInt($street, 'id'),
                $this->getString($street, 'name'),
                $buildings,
                $comment,
            );
        }

        return array_values($result);
    }

    /**
     * @param array<string, mixed> $data
     */
    private function getString(array $data, string $key): string
    {
        $value = $data[$key] ?? '';

        return is_string($value) || is_numeric($value) ? (string) $value : '';
    }

    /**
     * @param array<string, mixed> $data
     */
    private function getInt(array $data, string $key): int
    {
        $value = $data[$key] ?? 0;

        return is_numeric($value) ? (int) $value : 0;
    }

    private function castToString(mixed $value): string
    {
        return is_string($value) || is_numeric($value) ? (string) $value : '';
    }
}
