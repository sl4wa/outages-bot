<?php

declare(strict_types=1);

namespace App\Tests\Support;

use App\Application\Interface\Repository\UserRepositoryInterface;
use App\Domain\Entity\User;

final class TestUserRepository implements UserRepositoryInterface
{
    /** @var User[] */
    public array $all = [];
    /** @var User[] */
    public array $saved = [];
    /** @var int[] */
    public array $removed = [];

    public function findAll(): array { return $this->all; }
    public function find(int $chatId): ?User { return null; }
    public function save(User $user): void { $this->saved[] = $user; }
    public function remove(int $chatId): void { $this->removed[] = $chatId; }
}
