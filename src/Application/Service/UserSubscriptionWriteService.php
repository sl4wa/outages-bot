<?php

namespace App\Application\Service;

use App\Application\DTO\UserSubscriptionDTO;
use App\Application\Interface\Repository\UserRepositoryInterface;
use App\Domain\Entity\User;

readonly class UserSubscriptionWriteService
{
    public function __construct(private UserRepositoryInterface $userRepository) {}

    public function createOrUpdate(int $chatId, int $streetId, string $streetName, string $building): UserSubscriptionDTO
    {
        $user = new User(
            id: $chatId,
            streetId: $streetId,
            streetName: $streetName,
            building: $building,
            startDate: null,
            endDate: null,
            comment: ''
        );

        $this->userRepository->save($user);

        return new UserSubscriptionDTO(
            chatId: $user->id,
            streetId: $user->streetId,
            streetName: $user->streetName,
            building: $user->building,
        );
    }
}
