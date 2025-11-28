<?php

declare(strict_types=1);

namespace App\Domain\Service;

use App\Domain\Entity\Outage;
use App\Domain\Entity\User;
use App\Domain\ValueObject\OutageInfo;

final class OutageFinder
{
    /**
     * @param Outage[] $allOutages
     */
    public function findOutageForNotification(User $user, array $allOutages): ?Outage
    {
        foreach ($allOutages as $outage) {
            if (!$outage->affectsUserAddress($user->address)) {
                continue;
            }

            $outageInfo = new OutageInfo($outage->period, $outage->description);

            if ($user->isAlreadyNotifiedAbout($outageInfo)) {
                return null;
            }

            return $outage; // Found first matching outage
        }

        return null;
    }
}
