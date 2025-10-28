<?php

namespace App\Application\Notifier\Command;

use App\Application\Interface\Repository\UserRepositoryInterface;
use App\Domain\Entity\User;
use App\Domain\ValueObject\OutageData;

readonly class RecordUserNotificationCommandHandler
{
    public function __construct(
        private UserRepositoryInterface $userRepository,
    ) {}

    public function handle(User $user, OutageData $outageData): void
    {
        $updatedUser = $user->withNotifiedOutage($outageData);
        $this->userRepository->save($updatedUser);
    }
}
