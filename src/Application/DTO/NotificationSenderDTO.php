<?php

namespace App\Application\DTO;

use App\Domain\Entity\Outage;
use App\Domain\Entity\User;

readonly class NotificationSenderDTO
{
    public function __construct(
        public User $user,
        public Outage $outage
    ) {}
}
