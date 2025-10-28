<?php

namespace App\Application\Notifier\Service;

use App\Application\Notifier\Interface\Provider\OutageProviderInterface;
use App\Domain\Entity\Outage;
use App\Domain\ValueObject\Address;
use App\Domain\ValueObject\OutageData;

readonly class OutageFetchService
{
    public function __construct(
        private OutageProviderInterface $outageProvider,
    ) {}

    /**
     * @return Outage[]
     */
    public function handle(): array
    {
        $dtos = $this->outageProvider->fetchOutages();
        return array_map(
            fn($dto) => new Outage(
                new OutageData($dto->start, $dto->end, $dto->comment),
                new Address($dto->streetId, $dto->streetName, $dto->buildings, $dto->city)
            ),
            $dtos
        );
    }
}

