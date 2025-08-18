<?php

declare(strict_types=1);

namespace App\Tests\Support;

use App\Application\Interface\Provider\OutageProviderInterface;
use App\Domain\Entity\Outage;

final class TestOutageProvider implements OutageProviderInterface
{
    /** @var Outage[] */
    public array $outages = [];

    public function fetchOutages(): array
    {
        return $this->outages;
    }
}

