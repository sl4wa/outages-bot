<?php
namespace App\Domain\Entity;

use App\Domain\ValueObject\OutageInfo;
use App\Domain\ValueObject\UserAddress;

readonly class User
{
    public function __construct(
        public int $id,
        public UserAddress $address,
        public ?OutageInfo $outageInfo,
    ) {}

    public function withNotifiedOutage(Outage $outage): self
    {
        $outageInfo = new OutageInfo($outage->period, $outage->description);

        return new self(
            $this->id,
            $this->address,
            $outageInfo
        );
    }

    public function isAlreadyNotifiedAbout(OutageInfo $outageInfo): bool
    {
        if ($this->outageInfo === null) {
            return false;
        }

        return $this->outageInfo->equals($outageInfo);
    }
}

