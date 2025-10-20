<?php

declare(strict_types=1);

namespace App\Tests\Support;

use App\Application\DTO\NotificationSenderDTO;
use App\Application\Exception\NotificationSendException;
use App\Application\Interface\Service\NotificationSenderInterface;

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

