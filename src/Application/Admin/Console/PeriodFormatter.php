<?php

declare(strict_types=1);

namespace App\Application\Admin\Console;

use DateTimeImmutable;

final class PeriodFormatter
{
    public const string DATE_FORMAT = 'd.m.Y';

    public const string TIME_FORMAT = 'H:i';

    public const string DATETIME_FORMAT = self::DATE_FORMAT . ' ' . self::TIME_FORMAT;

    public static function format(DateTimeImmutable $start, DateTimeImmutable $end): string
    {
        if ($start->format('Y-m-d') === $end->format('Y-m-d')) {
            return $start->format(self::DATETIME_FORMAT) . ' - ' . $end->format(self::TIME_FORMAT);
        }

        return $start->format(self::DATETIME_FORMAT) . ' - ' . $end->format(self::DATETIME_FORMAT);
    }
}
