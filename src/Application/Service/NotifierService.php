<?php
namespace App\Application\Service;

use App\Application\DTO\NotificationSenderDTO;
use App\Application\Exception\NotificationSendException;
use App\Application\Interface\Repository\UserRepositoryInterface;
use App\Application\Interface\Service\NotificationSenderInterface;
use App\Domain\Service\OutageFinder;

class NotifierService
{
    public function __construct(
        private readonly OutageFetchService $outageFetchService,
        private readonly UserRepositoryInterface $userRepository,
        private readonly NotificationSenderInterface $notificationSender,
        private readonly OutageFinder $outageFinder,
    ) {}

    public function notify(): int
    {
        $outages = $this->outageFetchService->fetch();
        $users = $this->userRepository->findAll();

        foreach ($users as $user) {
            $outageToNotify = $this->outageFinder->findOutageForNotification($user, $outages);

            if ($outageToNotify !== null) {
                try {
                    $this->notificationSender->send(new NotificationSenderDTO($user, $outageToNotify));
                    $updatedUser = $user->withUpdatedOutage($outageToNotify);
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
