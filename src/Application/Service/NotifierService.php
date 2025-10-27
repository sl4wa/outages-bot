<?php
namespace App\Application\Service;

use App\Application\DTO\NotificationSenderDTO;
use App\Application\Exception\NotificationSendException;
use App\Application\Interface\Repository\UserRepositoryInterface;
use App\Application\Interface\Service\NotificationSenderInterface;
use App\Domain\Entity\Outage;
use App\Domain\Entity\User;
use App\Domain\Service\OutageFinder;
use App\Domain\ValueObject\OutageData;

readonly class NotifierService
{
    public function __construct(
        private UserRepositoryInterface $userRepository,
        private NotificationSenderInterface $notificationSender,
        private OutageFinder $outageFinder,
    ) {}

    /**
     * @param User[] $users
     * @param Outage[] $outages
     */
    public function notify(array $users, array $outages): int
    {
        foreach ($users as $user) {
            $outageToNotify = $this->outageFinder->findOutageForNotification($user, $outages);

            if ($outageToNotify !== null) {
                try {
                    $this->notificationSender->send(new NotificationSenderDTO(
                        userId: $user->id,
                        city: $outageToNotify->address->city,
                        streetName: $outageToNotify->address->streetName,
                        buildings: $outageToNotify->address->buildings,
                        start: $outageToNotify->data->startDate,
                        end: $outageToNotify->data->endDate,
                        comment: $outageToNotify->data->comment,
                    ));
                    $updatedUser = $user->withNotifiedOutage($outageToNotify->data);
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
