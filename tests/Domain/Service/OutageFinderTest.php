<?php

declare(strict_types=1);

namespace App\Tests\Domain\Service;

use App\Domain\Entity\Outage;
use App\Domain\Entity\User;
use App\Domain\Service\OutageFinder;
use App\Domain\ValueObject\Address;
use App\Domain\ValueObject\OutageData;
use PHPUnit\Framework\TestCase;

final class OutageFinderTest extends TestCase
{
    private const TEST_BUILDINGS = '271, 273, 273-А, 275, 277, 279, 281, 281-А, 282, 283, 283-А, '
        . '284, 284-А, 285, 285-А, 287, 289, 289-А, 290-А, 291, 291(0083), '
        . '293, 295, 297, 297-А, 297-Б, 308, 313, 316, 316-А, 318, 318-А, '
        . '320, 322, 324, 326, 328, 328-А, 330, 332, 334, 336, 338, 340-А, '
        . '342, 346, 348-А, 350, 350,А, 350-В, 358, 358-А, 360-В';

    private OutageFinder $finder;
    private Outage $outage1;
    private Outage $outage2;

    protected function setUp(): void
    {
        $this->finder = new OutageFinder();

        $outageData1 = new OutageData(
            startDate: new \DateTimeImmutable('2024-11-28T06:47:00+00:00'),
            endDate: new \DateTimeImmutable('2024-11-28T10:00:00+00:00'),
            comment: 'Застосування ГПВ'
        );

        $address1 = new Address(
            streetId: 12783,
            streetName: 'Шевченка Т.',
            buildings: array_map('trim', explode(',', self::TEST_BUILDINGS)),
            city: 'Львів'
        );

        $this->outage1 = new Outage(
            data: $outageData1,
            address: $address1,
        );

        $outageData2 = new OutageData(
            startDate: new \DateTimeImmutable('2024-11-28T06:47:00+00:00'),
            endDate: new \DateTimeImmutable('2024-11-28T10:00:00+00:00'),
            comment: 'Застосування ГПВ'
        );

        $address2 = new Address(
            streetId: 6458,
            streetName: 'Хмельницького Б.',
            buildings: ['294'],
            city: 'Львів'
        );

        $this->outage2 = new Outage(
            data: $outageData2,
            address: $address2,
        );
    }

    public function testFindsMatchingOutageForUser(): void
    {
        $user1 = new User(
            id: 1,
            address: new Address(streetId: 12783, streetName: 'Шевченка Т.', buildings: ['271']),
            lastNotifiedOutage: null
        );
        $user2 = new User(
            id: 2,
            address: new Address(streetId: 12783, streetName: 'Шевченка Т.', buildings: ['279']),
            lastNotifiedOutage: null
        );
        $user3 = new User(
            id: 3,
            address: new Address(streetId: 6458, streetName: 'Хмельницького Б.', buildings: ['294']),
            lastNotifiedOutage: null
        );

        $allOutages = [$this->outage1, $this->outage2];

        $result1 = $this->finder->findOutageForNotification($user1, $allOutages);
        self::assertSame($this->outage1, $result1);

        $result2 = $this->finder->findOutageForNotification($user2, $allOutages);
        self::assertSame($this->outage1, $result2);

        $result3 = $this->finder->findOutageForNotification($user3, $allOutages);
        self::assertSame($this->outage2, $result3);
    }

    public function testNoMatchingOutageReturnsNull(): void
    {
        $user = new User(
            id: 1,
            address: new Address(streetId: 13961, streetName: 'Залізнична', buildings: ['16']),
            lastNotifiedOutage: null
        );
        $result = $this->finder->findOutageForNotification($user, [$this->outage1, $this->outage2]);
        self::assertNull($result);
    }

    public function testAlreadyAwareUserReturnsNull(): void
    {
        $already = new User(
            id: 10,
            address: new Address(streetId: 12783, streetName: 'Шевченка Т.', buildings: ['271']),
            lastNotifiedOutage: $this->outage1->data
        );

        $result = $this->finder->findOutageForNotification($already, [$this->outage1, $this->outage2]);
        self::assertNull($result);
    }
}
