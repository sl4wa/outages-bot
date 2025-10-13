<?php

namespace App\Application\Service;

use App\Application\Factory\OutageFactory;
use App\Application\Interface\Provider\OutageProviderInterface;
use App\Domain\Entity\Outage;

readonly class OutageFetchService
{
    public function __construct(
        private OutageProviderInterface $outageProvider,
        private OutageFactory $outageFactory,
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

