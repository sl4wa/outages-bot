<?php

declare(strict_types=1);

namespace App\Infrastructure\Provider;

use App\Application\Notifier\DTO\OutageDTO;
use App\Application\Notifier\Interface\Provider\OutageProviderInterface;
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
        $items = $data['hydra:member'] ?? [];

        return array_map(function (array $row): OutageDTO {
            $comment = (string) ($row['koment'] ?? '');
            $comment = preg_replace('/[\r\n]+/', ' ', $comment);
            $comment = trim($comment);

            $buildingsRaw = $row['buildingNames'] ?? '';
            $buildings = is_array($buildingsRaw)
                ? array_map('trim', $buildingsRaw)
                : array_filter(array_map('trim', explode(',', (string) $buildingsRaw)));

            return new OutageDTO(
                (int) ($row['id'] ?? 0),
                new DateTimeImmutable($row['dateEvent']),
                new DateTimeImmutable($row['datePlanIn']),
                (string) ($row['city']['name'] ?? ''),
                (int) ($row['street']['id'] ?? 0),
                (string) ($row['street']['name'] ?? ''),
                $buildings,
                $comment,
            );
        }, $items);
    }
}
