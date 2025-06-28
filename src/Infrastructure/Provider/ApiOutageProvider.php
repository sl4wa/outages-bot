<?php
namespace App\Infrastructure\Provider;

use App\Application\Interface\Provider\OutageProviderInterface;
use App\Domain\Entity\Outage;
use Symfony\Contracts\HttpClient\HttpClientInterface;

class ApiOutageProvider implements OutageProviderInterface
{
    private const API_URL = 'https://power-api.loe.lviv.ua/api/pw_accidents?pagination=false&otg.id=28&city.id=693';

    public function __construct(private readonly HttpClientInterface $httpClient) {}

    /**
     * @return Outage[]
     */
    public function fetchOutages(): array
    {
        $response = $this->httpClient->request('GET', self::API_URL);
        if ($response->getStatusCode() !== 200) {
            return [];
        }
        $data = $response->toArray();
        return array_map([Outage::class, 'fromArray'], $data['hydra:member'] ?? []);
    }
}
