<?php

declare(strict_types=1);

namespace App\Application\Notifier\Interface\Provider;

use App\Application\Notifier\DTO\OutageDTO;

interface OutageProviderInterface
{
    /**
     * @return OutageDTO[]
     */
    public function fetchOutages(): array;
}
