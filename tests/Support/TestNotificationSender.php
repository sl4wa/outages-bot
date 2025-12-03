<?php

declare(strict_types=1);

namespace App\Tests\Support;

use App\Application\Interface\NotificationSenderInterface;
use App\Application\Notifier\DTO\NotificationSenderDTO;
use App\Application\Notifier\Exception\NotificationSendException;

final class TestNotificationSender implements NotificationSenderInterface
{
    /** @var NotificationSenderDTO[] */
    public array $sent = [];

    public ?int $blockUserId = null;

    public function send(NotificationSenderDTO $dto): void
    {
        if ($this->blockUserId !== null && $dto->userId === $this->blockUserId) {
            throw new NotificationSendException($dto->userId, 'Forbidden', 403);
        }
        $this->sent[] = $dto;
    }
}
