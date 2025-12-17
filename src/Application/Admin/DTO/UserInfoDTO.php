<?php

declare(strict_types=1);

namespace App\Application\Admin\DTO;

final readonly class UserInfoDTO
{
    public function __construct(
        public int $chatId,
        public ?string $username,
        public ?string $firstName,
        public ?string $lastName,
    ) {
    }
}
