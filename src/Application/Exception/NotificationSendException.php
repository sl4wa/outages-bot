<?php

namespace App\Application\Exception;

use Throwable;

class NotificationSendException extends \RuntimeException
{
    public function __construct(
        public readonly int $userId,
        string $message = '',
        int $code = 0,
        ?Throwable $previous = null
    ) {
        parent::__construct($message, $code, $previous);
    }

    public function isBlocked(): bool
    {
        return $this->code === 403 || str_contains(strtolower($this->message), 'forbidden');
    }
}
