<?php

declare(strict_types=1);

namespace App\Application\Notifier\Service;

use App\Application\Notifier\Interface\Provider\OutageProviderInterface;
use App\Domain\Entity\Outage;
use App\Domain\ValueObject\OutageAddress;
use App\Domain\ValueObject\OutageDescription;
use App\Domain\ValueObject\OutagePeriod;

final readonly class OutageFetchService
{
    public function __construct(
        private OutageProviderInterface $outageProvider,
    ) {
    }

    /**
     * @return Outage[]
     */
    public function handle(): array
    {
        $dtos = $this->outageProvider->fetchOutages();

        return array_map(
            fn ($dto) => new Outage(
                $dto->id,
                new OutagePeriod($dto->start, $dto->end),
                new OutageAddress($dto->streetId, $dto->streetName, $dto->buildings, $dto->city),
                new OutageDescription($dto->comment)
            ),
            $dtos
        );
    }
}
