<?php

namespace App\Application\Interface\Service;

use App\Application\DTO\NotificationSenderDTO;
use App\Application\Exception\NotificationSendException;

interface NotificationSenderInterface
{
    /**
     * @throws NotificationSendException
     */
    public function send(NotificationSenderDTO $dto): void;
}
