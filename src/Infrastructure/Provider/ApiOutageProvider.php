<?php

namespace App\Infrastructure\Provider;

use App\Application\DTO\OutageDTO;
use App\Application\Interface\Provider\OutageProviderInterface;
use Symfony\Contracts\HttpClient\HttpClientInterface;

class ApiOutageProvider implements OutageProviderInterface
{
    private const API_URL = 'https://power-api.loe.lviv.ua/api/pw_accidents?pagination=false&otg.id=28&city.id=693';

    public function __construct(private readonly HttpClientInterface $httpClient) {}

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
            $comment = (string)($row['koment'] ?? '');
            $comment = preg_replace('/[\r\n]+/', ' ', $comment);
            $comment = trim($comment);

            $buildings = $row['buildingNames'] ?? '';
            $buildingNames = is_array($buildings)
                ? array_map('trim', $buildings)
                : array_filter(array_map('trim', explode(',', (string)$buildings)));

            return new OutageDTO(
                new \DateTimeImmutable($row['dateEvent']),
                new \DateTimeImmutable($row['datePlanIn']),
                (string)($row['city']['name'] ?? ''),
                (int)($row['street']['id'] ?? 0),
                (string)($row['street']['name'] ?? ''),
                $buildingNames,
                $comment,
            );
        }, $items);
    }
}
