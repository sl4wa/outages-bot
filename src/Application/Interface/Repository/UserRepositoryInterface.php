<?php


namespace App\Application\Interface\Repository;

use App\Domain\Entity\User;

interface UserRepositoryInterface
{
    public function findAll(): array;
    public function find(int $chatId): ?User;
    public function save(User $user): void;
    public function remove(int $chatId): void;
}
