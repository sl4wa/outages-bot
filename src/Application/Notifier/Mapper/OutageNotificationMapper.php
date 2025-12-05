<?php

declare(strict_types=1);

namespace App\Application\Notifier\Mapper;

use App\Application\Notifier\DTO\NotificationSenderDTO;
use App\Domain\Entity\Outage;

final class OutageNotificationMapper
{
    public function toNotificationDTO(Outage $outage, int $userId): NotificationSenderDTO
    {
        return new NotificationSenderDTO(
            userId: $userId,
            city: $outage->address->city ?? '',
            streetName: $outage->address->streetName,
            buildings: array_values($outage->address->buildings),
            start: $outage->period->startDate,
            end: $outage->period->endDate,
            comment: $outage->description->value,
        );
    }
}
