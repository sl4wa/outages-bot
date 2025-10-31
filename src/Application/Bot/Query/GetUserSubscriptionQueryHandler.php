<?php

namespace App\Application\Bot\Query;

use App\Application\Bot\DTO\UserSubscriptionDTO;
use App\Application\Interface\Repository\UserRepositoryInterface;

readonly class GetUserSubscriptionQueryHandler
{
    public function __construct(private UserRepositoryInterface $userRepository) {}

    public function handle(int $chatId): ?UserSubscriptionDTO
    {
        $user = $this->userRepository->find($chatId);
        if ($user === null) {
            return null;
        }

        return new UserSubscriptionDTO(
            chatId: $user->id,
            streetId: $user->address->streetId,
            streetName: $user->address->streetName,
            building: $user->address->building,
        );
    }
}

