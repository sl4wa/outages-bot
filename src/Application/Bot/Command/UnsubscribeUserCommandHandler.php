<?php

namespace App\Application\Bot\Command;

use App\Application\Interface\Repository\UserRepositoryInterface;

readonly class UnsubscribeUserCommandHandler
{
    public function __construct(private UserRepositoryInterface $userRepository) {}

    public function handle(int $chatId): bool
    {
        $user = $this->userRepository->find($chatId);
        if ($user === null) {
            return false;
        }

        $this->userRepository->remove($chatId);
        return true;
    }
}

