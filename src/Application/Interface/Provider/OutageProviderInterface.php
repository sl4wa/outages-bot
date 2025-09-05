<?php

declare(strict_types=1);

namespace App\Application\Interface\Provider;

use App\Application\DTO\OutageDTO;

interface OutageProviderInterface
{
    /**
     * @return OutageDTO[]
     */
    public function fetchOutages(): array;
}
