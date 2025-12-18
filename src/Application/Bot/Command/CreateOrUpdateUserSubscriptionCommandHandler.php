<?php

declare(strict_types=1);

namespace App\Application\Bot\Command;

use App\Application\Bot\DTO\UserSubscriptionDTO;
use App\Domain\Entity\User;
use App\Domain\Repository\UserRepositoryInterface;
use App\Domain\ValueObject\UserAddress;

final readonly class CreateOrUpdateUserSubscriptionCommandHandler
{
    public function __construct(private UserRepositoryInterface $userRepository)
    {
    }

    public function handle(int $chatId, int $streetId, string $streetName, string $building): UserSubscriptionDTO
    {
        $address = new UserAddress($streetId, $streetName, $building);

        $user = new User(
            id: $chatId,
            address: $address,
            outageInfo: null
        );

        $this->userRepository->save($user);

        return new UserSubscriptionDTO(
            chatId: $user->id,
            streetId: $user->address->streetId,
            streetName: $user->address->streetName,
            building: $user->address->building,
        );
    }
}
