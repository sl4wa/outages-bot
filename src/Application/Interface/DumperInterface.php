<?php

declare(strict_types=1);

namespace App\Application\Interface;

interface DumperInterface
{
    public function dump(mixed $data, string $filename): void;
}
