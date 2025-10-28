<?php

namespace App\Application\Notifier\Command;

use App\Application\Interface\Repository\UserRepositoryInterface;

readonly class RemoveUserCommandHandler
{
    public function __construct(
        private UserRepositoryInterface $userRepository,
    ) {}

    public function handle(int $userId): void
    {
        $this->userRepository->remove($userId);
    }
}
