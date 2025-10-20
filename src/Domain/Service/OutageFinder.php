<?php

namespace App\Domain\Service;

use App\Domain\Entity\Outage;
use App\Domain\Entity\User;

final class OutageFinder
{
    /**
     * @param User $user
     * @param Outage[] $allOutages
     * @return ?Outage
     */
    public function findOutageForNotification(User $user, array $allOutages): ?Outage
    {
        foreach ($allOutages as $outage) {
            if (!$outage->address->covers($user->address)) {
                continue;
            }

            if ($outage->isIdenticalPeriodAndComment($user)) {
                return null; // User already aware
            }

            return $outage; // Found first matching outage
        }

        return null;
    }
}
