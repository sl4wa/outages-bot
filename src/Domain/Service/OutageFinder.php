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
        $matchingOutages = array_filter(
            $allOutages,
            static fn(Outage $outage) => $outage->matchesUser($user)
        );

        if (empty($matchingOutages)) {
            return null;
        }

        foreach ($matchingOutages as $outage) {
            if ($outage->isIdenticalPeriodAndComment($user)) {
                return null; // The user is already aware of an active outage.
            }
        }

        // A new or changed outage situation exists. Notify about the first one found.
        return reset($matchingOutages);
    }
}
