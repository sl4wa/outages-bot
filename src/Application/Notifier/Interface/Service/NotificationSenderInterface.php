<?php

declare(strict_types=1);

namespace App\Application\Notifier\Interface\Service;

use App\Application\Notifier\DTO\NotificationSenderDTO;
use App\Application\Notifier\Exception\NotificationSendException;

interface NotificationSenderInterface
{
    /**
     * @throws NotificationSendException
     */
    public function send(NotificationSenderDTO $dto): void;
}
