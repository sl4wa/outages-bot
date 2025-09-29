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
        $potentialOutageToNotify = null;

        foreach ($allOutages as $outage) {
            if (!$outage->matchesUser($user)) {
                continue;
            }

            if ($outage->isIdenticalPeriodAndComment($user)) {
                // The user is already aware of an active outage.
                return null;
            }

            if ($potentialOutageToNotify === null) {
                $potentialOutageToNotify = $outage;
            }
        }

        return $potentialOutageToNotify;
    }
}
