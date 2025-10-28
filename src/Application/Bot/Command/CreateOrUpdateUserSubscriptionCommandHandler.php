<?php

namespace App\Application\Bot\Command;

use App\Application\Bot\DTO\UserSubscriptionDTO;
use App\Application\Interface\Repository\UserRepositoryInterface;
use App\Domain\Entity\User;
use App\Domain\ValueObject\Address;

readonly class CreateOrUpdateUserSubscriptionCommandHandler
{
    public function __construct(private UserRepositoryInterface $userRepository) {}

    public function handle(int $chatId, int $streetId, string $streetName, string $building): UserSubscriptionDTO
    {
        $address = new Address($streetId, $streetName, [$building]);

        $user = new User(
            id: $chatId,
            address: $address,
            lastNotifiedOutage: null
        );

        $this->userRepository->save($user);

        return new UserSubscriptionDTO(
            chatId: $user->id,
            streetId: $user->address->streetId,
            streetName: $user->address->streetName,
            building: $user->address->getSingleBuilding(),
        );
    }
}
