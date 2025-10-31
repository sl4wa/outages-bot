<?php

namespace App\Application\Notifier\Command;

use App\Application\Notifier\Exception\NotificationSendException;
use App\Application\Notifier\Interface\Service\NotificationSenderInterface;
use App\Application\Notifier\Mapper\OutageNotificationMapper;
use App\Domain\Entity\Outage;
use App\Domain\Entity\User;

readonly class SendNotificationCommandHandler
{
    public function __construct(
        private NotificationSenderInterface $notificationSender,
        private OutageNotificationMapper $mapper,
    ) {}

    /**
     * @throws NotificationSendException
     */
    public function handle(User $user, Outage $outage): void
    {
        $this->notificationSender->send(
            $this->mapper->toNotificationDTO($outage, $user->id)
        );
    }
}
