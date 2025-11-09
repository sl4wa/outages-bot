<?php

declare(strict_types=1);

namespace App\Infrastructure\Telegram\Bot;

use App\Application\Bot\DTO\UserInfoDTO;
use App\Application\Bot\Interface\TelegramUserInfoProviderInterface;
use SergiX44\Nutgram\Nutgram;

final readonly class TelegramUserInfoProvider implements TelegramUserInfoProviderInterface
{
    public function __construct(
        private Nutgram $nutgram
    ) {
    }

    public function getUserInfo(int $chatId): UserInfoDTO
    {
        try {
            $chat = $this->nutgram->getChat($chatId);

            return new UserInfoDTO(
                chatId: $chat->id,
                username: $chat->username,
                firstName: $chat->first_name,
                lastName: $chat->last_name,
            );
        } catch (\Throwable $e) {
            throw new \RuntimeException(
                "Failed to get user info for chat {$chatId}: {$e->getMessage()}",
                0,
                $e
            );
        }
    }
}
