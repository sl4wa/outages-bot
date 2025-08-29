<?php

namespace App\Application\Service;

use App\Application\DTO\UserSubscriptionDTO;
use App\Application\Interface\Repository\UserRepositoryInterface;

class UserSubscriptionQueryService
{
    public function __construct(private readonly UserRepositoryInterface $userRepository) {}

    public function get(int $chatId): ?UserSubscriptionDTO
    {
        $user = $this->userRepository->find($chatId);
        if ($user === null) {
            return null;
        }

        return new UserSubscriptionDTO(
            chatId: $user->id,
            streetId: $user->streetId,
            streetName: $user->streetName,
            building: $user->building,
        );
    }
}

