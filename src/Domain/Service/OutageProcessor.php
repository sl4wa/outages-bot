<?php

namespace App\Domain\Service;

use App\Domain\Entity\Outage;
use App\Domain\Entity\User;

class OutageProcessor
{
    /**
     * @param Outage $outage
     * @param User[] $usersToBeChecked
     * @return User[]
     */
    public function process(Outage $outage, array $usersToBeChecked): array
    {
        $usersToBeNotified = [];

        foreach ($usersToBeChecked as $user) {
            if ($outage->matchesUser($user) && !$outage->isIdenticalPeriodAndComment($user)) {
                $usersToBeNotified[] = $user;
            }
        }
        return $usersToBeNotified;
    }
}
