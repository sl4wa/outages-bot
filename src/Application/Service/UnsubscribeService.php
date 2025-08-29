<?php

namespace App\Application\Service;

use App\Application\Interface\Repository\UserRepositoryInterface;

class UnsubscribeService
{
    public function __construct(private readonly UserRepositoryInterface $userRepository) {}

    /**
     * Removes subscription if exists.
     *
     * @return bool True if removed, false if not found
     */
    public function unsubscribe(int $chatId): bool
    {
        $user = $this->userRepository->find($chatId);
        if ($user === null) {
            return false;
        }

        $this->userRepository->remove($chatId);
        return true;
    }
}

