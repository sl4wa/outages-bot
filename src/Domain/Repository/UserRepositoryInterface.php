<?php

declare(strict_types=1);

namespace App\Domain\Repository;

use App\Domain\Entity\User;

interface UserRepositoryInterface
{
    /**
     * @return User[]
     */
    public function findAll(): array;

    public function find(int $chatId): ?User;

    public function save(User $user): void;

    public function remove(int $chatId): bool;
}
