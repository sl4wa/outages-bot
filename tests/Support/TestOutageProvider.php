<?php

declare(strict_types=1);

namespace App\Tests\Support;

use App\Application\Notifier\DTO\OutageDTO;
use App\Application\Notifier\Interface\Provider\OutageProviderInterface;

final class TestOutageProvider implements OutageProviderInterface
{
    /** @var OutageDTO[] */
    public array $outages = [];

    public function fetchOutages(): array
    {
        return $this->outages;
    }
}
