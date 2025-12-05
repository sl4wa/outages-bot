<?php

declare(strict_types=1);

namespace App\Application\Interface;

use App\Application\Bot\DTO\UserInfoDTO;

interface TelegramUserInfoProviderInterface
{
    public function getUserInfo(int $chatId): UserInfoDTO;
}
