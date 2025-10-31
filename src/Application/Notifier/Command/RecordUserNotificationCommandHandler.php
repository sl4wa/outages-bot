<?php

namespace App\Application\Notifier\Command;

use App\Application\Interface\Repository\UserRepositoryInterface;
use App\Domain\Entity\User;
use App\Domain\ValueObject\OutageInfo;

readonly class RecordUserNotificationCommandHandler
{
    public function __construct(
        private UserRepositoryInterface $userRepository,
    ) {}

    public function handle(User $user, OutageInfo $outageInfo): void
    {
        $updatedUser = $user->withNotifiedOutage($outageInfo);
        $this->userRepository->save($updatedUser);
    }
}
