<?php

declare(strict_types=1);

namespace App\Application\Interface\Provider;

use App\Domain\Entity\Outage;

interface OutageProviderInterface
{
    /**
     * @return Outage[]
     */
    public function fetchOutages(): array;
}
