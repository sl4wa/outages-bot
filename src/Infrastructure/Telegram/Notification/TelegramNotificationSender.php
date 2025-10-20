<?php

namespace App\Infrastructure\Telegram\Notification;

use App\Application\DTO\NotificationSenderDTO;
use App\Application\Exception\NotificationSendException;
use App\Application\Interface\Service\NotificationSenderInterface;
use SergiX44\Nutgram\Nutgram;

readonly class TelegramNotificationSender implements NotificationSenderInterface
{
    public function __construct(
        private Nutgram $bot,
        private TelegramNotificationFormatter $formatter,
    ) {}

    public function send(NotificationSenderDTO $dto): void
    {
        try {
            $this->bot->sendMessage(
                text: $this->formatter->format($dto),
                chat_id: $dto->userId,
                parse_mode: 'HTML'
            );
        } catch (\Throwable $e) {
            throw new NotificationSendException(
                $dto->userId,
                $e->getMessage(),
                0,
                $e
            );
        }
    }
}
