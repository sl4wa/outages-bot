<?php
namespace App\Application\Service;

use App\Application\DTO\NotificationSenderDTO;
use App\Application\Exception\NotificationSendException;
use App\Application\Interface\Provider\OutageProviderInterface;
use App\Application\Interface\Repository\UserRepositoryInterface;
use App\Application\Interface\Service\NotificationSenderInterface;
use App\Domain\Service\OutageProcessor;

class NotifierService
{
    public function __construct(
        private readonly OutageProviderInterface $outageProvider,
        private readonly UserRepositoryInterface $userRepository,
        private readonly NotificationSenderInterface $notificationSender,
        private readonly OutageProcessor $outageProcessor,
    ) {}

    public function notify(): int
    {
        $outages = $this->outageProvider->fetchOutages();
        $usersToBeChecked = $this->userRepository->findAll();

        $notifiedUserIds = [];

        foreach ($outages as $outage) {
            $usersToBeNotified = $this->outageProcessor->process($outage, $usersToBeChecked);

            foreach ($usersToBeNotified as $user) {
                if (in_array($user->id, $notifiedUserIds, true)) {
                    continue;
                }

                try {
                    $this->notificationSender->send(new NotificationSenderDTO($user, $outage));
                    $updatedUser = $user->withUpdatedOutage($outage);
                    $this->userRepository->save($updatedUser);
                    $notifiedUserIds[] = $user->id;
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
