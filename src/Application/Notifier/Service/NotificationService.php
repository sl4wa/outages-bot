<?php
namespace App\Application\Notifier\Service;

use App\Application\Notifier\Command\RecordUserNotificationCommandHandler;
use App\Application\Notifier\Command\RemoveUserCommandHandler;
use App\Application\Notifier\Command\SendNotificationCommandHandler;
use App\Application\Notifier\Exception\NotificationSendException;
use App\Domain\Entity\Outage;
use App\Domain\Entity\User;
use App\Domain\Service\OutageFinder;

readonly class NotificationService
{
    public function __construct(
        private SendNotificationCommandHandler $sendNotificationCommandHandler,
        private RecordUserNotificationCommandHandler $recordUserNotificationCommandHandler,
        private RemoveUserCommandHandler $removeUserCommandHandler,
        private OutageFinder $outageFinder,
    ) {}

    /**
     * @param User[] $users
     * @param Outage[] $outages
     */
    public function handle(array $users, array $outages): int
    {
        foreach ($users as $user) {
            $outageToNotify = $this->outageFinder->findOutageForNotification($user, $outages);

            if ($outageToNotify !== null) {
                try {
                    $this->sendNotificationCommandHandler->handle($user, $outageToNotify);
                    $this->recordUserNotificationCommandHandler->handle($user, $outageToNotify->data);
                } catch (NotificationSendException $e) {
                    if ($e->isBlocked()) {
                        $this->removeUserCommandHandler->handle($e->userId);
                    }
                }
            }
        }

        return count($outages);
    }
}
