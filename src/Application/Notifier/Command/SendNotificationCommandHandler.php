<?php

namespace App\Application\Notifier\Command;

use App\Application\Notifier\DTO\NotificationSenderDTO;
use App\Application\Notifier\Exception\NotificationSendException;
use App\Application\Notifier\Interface\Service\NotificationSenderInterface;
use App\Domain\Entity\Outage;
use App\Domain\Entity\User;

readonly class SendNotificationCommandHandler
{
    public function __construct(
        private NotificationSenderInterface $notificationSender,
    ) {}

    /**
     * @throws NotificationSendException
     */
    public function handle(User $user, Outage $outage): void
    {
        $this->notificationSender->send(new NotificationSenderDTO(
            userId: $user->id,
            city: $outage->address->city,
            streetName: $outage->address->streetName,
            buildings: $outage->address->buildings,
            start: $outage->data->startDate,
            end: $outage->data->endDate,
            comment: $outage->data->comment,
        ));
    }
}
