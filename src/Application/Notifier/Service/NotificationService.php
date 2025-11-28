<?php

declare(strict_types=1);

namespace App\Application\Notifier\Service;

use App\Application\Interface\Repository\UserRepositoryInterface;
use App\Application\Notifier\Exception\NotificationSendException;
use App\Application\Notifier\Interface\Service\NotificationSenderInterface;
use App\Application\Notifier\Mapper\OutageNotificationMapper;
use App\Domain\Entity\Outage;
use App\Domain\Service\OutageFinder;

final readonly class NotificationService
{
    public function __construct(
        private NotificationSenderInterface $notificationSender,
        private OutageNotificationMapper $mapper,
        private UserRepositoryInterface $userRepository,
        private OutageFinder $outageFinder,
    ) {
    }

    /**
     * @param Outage[] $outages
     */
    public function handle(array $outages): int
    {
        $users = $this->userRepository->findAll();

        foreach ($users as $user) {
            $outageToNotify = $this->outageFinder->findOutageForNotification($user, $outages);

            if ($outageToNotify !== null) {
                try {
                    $this->notificationSender->send(
                        $this->mapper->toNotificationDTO($outageToNotify, $user->id)
                    );
                    $updatedUser = $user->withNotifiedOutage($outageToNotify);
                    $this->userRepository->save($updatedUser);
                } catch (NotificationSendException $e) {
                    if ($e->isBlocked()) {
                        $this->userRepository->remove($e->userId);
                    }
                }
            }
        }

        return count($outages);
    }
}
