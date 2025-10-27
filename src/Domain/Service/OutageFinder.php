<?php

namespace App\Domain\Service;

use App\Domain\Entity\Outage;
use App\Domain\Entity\User;
use App\Domain\ValueObject\OutageData;

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

            if ($user->wasAlreadyNotifiedAbout($outage->data)) {
                return null; // User already aware
            }

            return $outage; // Found first matching outage
        }

        return null;
    }
}
