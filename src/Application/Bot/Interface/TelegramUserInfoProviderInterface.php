<?php

declare(strict_types=1);

namespace App\Application\Bot\Interface;

use App\Application\Bot\DTO\UserInfoDTO;
use RuntimeException;

interface TelegramUserInfoProviderInterface
{
    /**
     * Get user information from Telegram by chat ID.
     *
     * @throws RuntimeException if user info cannot be retrieved
     */
    public function getUserInfo(int $chatId): UserInfoDTO;
}
