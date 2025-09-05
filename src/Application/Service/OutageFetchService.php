<?php

namespace App\Application\Service;

use App\Application\Factory\OutageFactory;
use App\Application\Interface\Provider\OutageProviderInterface;
use App\Domain\Entity\Outage;

class OutageFetchService
{
    public function __construct(
        private readonly OutageProviderInterface $outageProvider,
        private readonly OutageFactory $outageFactory,
    ) {}

    /**
     * @return Outage[]
     */
    public function fetch(): array
    {
        $dtos = $this->outageProvider->fetchOutages();
        return array_map(fn($dto) => $this->outageFactory->createFromDTO($dto), $dtos);
    }
}

