<?php
namespace App\Domain\Event;

use App\Domain\Entity\Outage;
use App\Domain\Entity\User;

final class OutageNotificationCreated
{
    /**
     * @param User[] $usersToBeNotified
     */
    public function __construct(
        public readonly Outage $outage,
        public readonly array $usersToBeNotified
    ) {}
}
