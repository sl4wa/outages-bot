<?php

declare(strict_types=1);

namespace App\Tests\Application\Service;

use App\Application\Notifier\DTO\OutageDTO;
use App\Application\Notifier\Service\NotificationService;
use App\Domain\Entity\Outage;
use App\Domain\Entity\User;
use App\Domain\ValueObject\OutageAddress;
use App\Domain\ValueObject\OutageDescription;
use App\Domain\ValueObject\OutagePeriod;
use App\Domain\ValueObject\UserAddress;
use App\Tests\Support\TestNotificationSender;
use App\Tests\Support\TestUserRepository;
use Symfony\Bundle\FrameworkBundle\Test\KernelTestCase;

final class NotifierServiceTest extends KernelTestCase
{
    private const TEST_BUILDINGS = '271, 273, 273-А, 275, 277, 279, 281, 281-А, 282, 283, 283-А, '
        . '284, 284-А, 285, 285-А, 287, 289, 289-А, 290-А, 291, 291(0083), '
        . '293, 295, 297, 297-А, 297-Б, 308, 313, 316, 316-А, 318, 318-А, '
        . '320, 322, 324, 326, 328, 328-А, 330, 332, 334, 336, 338, 340-А, '
        . '342, 346, 348-А, 350, 350,А, 350-В, 358, 358-А, 360-В';

    protected static function getKernelClass(): string
    {
        return \App\Kernel::class;
    }

    private NotificationService $notifier;
    private TestNotificationSender $sender;
    private TestUserRepository $userRepo;

    protected function setUp(): void
    {
        self::bootKernel();

        $container = self::getContainer();

        // Fetch services from the test container (DI)
        $this->sender = $container->get(TestNotificationSender::class);
        $this->userRepo = $container->get(TestUserRepository::class);
        $this->notifier = $container->get(NotificationService::class);

        // Reset state between tests (shared services)
        $this->sender->sent = [];
        $this->sender->blockUserId = null;
        $this->userRepo->saved = [];
        $this->userRepo->removed = [];
    }

    public function testNotificationSentAndUserSaved(): void
    {
        $outageDto = $this->createOutage('Застосування ГПВ');
        $outage = new Outage(
            $outageDto->id,
            new OutagePeriod($outageDto->start, $outageDto->end),
            new OutageAddress($outageDto->streetId, $outageDto->streetName, $outageDto->buildings, $outageDto->city),
            new OutageDescription($outageDto->comment)
        );
        $user = new User(
            id: 100,
            address: new UserAddress(streetId: 12783, streetName: 'Шевченка Т.', building: '271'),
            outageInfo: null
        );

        $this->notifier->handle([$user], [$outage]);

        self::assertCount(1, $this->sender->sent); // one notification
        self::assertEquals(100, $this->sender->sent[0]->userId);
        self::assertCount(1, $this->userRepo->saved);
        self::assertEquals('Застосування ГПВ', $this->sender->sent[0]->comment);
    }

    public function testSubscriptionRemovedForBlockedUser(): void
    {
        $this->sender->blockUserId = 101; // simulate Forbidden

        $outageDto = $this->createOutage('Застосування ГПВ');
        $outage = new Outage(
            $outageDto->id,
            new OutagePeriod($outageDto->start, $outageDto->end),
            new OutageAddress($outageDto->streetId, $outageDto->streetName, $outageDto->buildings, $outageDto->city),
            new OutageDescription($outageDto->comment)
        );
        $user = new User(
            id: 101,
            address: new UserAddress(streetId: 12783, streetName: 'Шевченка Т.', building: '271'),
            outageInfo: null
        );

        $this->notifier->handle([$user], [$outage]);

        self::assertSame([101], $this->userRepo->removed);
        self::assertCount(0, $this->sender->sent);
        self::assertCount(0, $this->userRepo->saved);
    }

    public function testNoRelevantOutageProducesNoNotification(): void
    {
        $outageDto = $this->createOutage('Застосування ГПВ');
        $outage = new Outage(
            $outageDto->id,
            new OutagePeriod($outageDto->start, $outageDto->end),
            new OutageAddress($outageDto->streetId, $outageDto->streetName, $outageDto->buildings, $outageDto->city),
            new OutageDescription($outageDto->comment)
        );
        $user = new User(
            id: 102,
            address: new UserAddress(streetId: 99999, streetName: 'Nonexistent Street', building: '1'),
            outageInfo: null
        );

        $this->notifier->handle([$user], [$outage]);

        self::assertCount(0, $this->sender->sent);
        self::assertCount(0, $this->userRepo->saved);
        self::assertCount(0, $this->userRepo->removed);
    }

    public function testMultipleOutagesForSameBuildingNotifiesOnlyOnce(): void
    {
        $user = new User(
            id: 103,
            address: new UserAddress(streetId: 12783, streetName: 'Шевченка Т.', building: '271'),
            outageInfo: null
        );
        $outageDtoA = $this->createOutage('Outage A');
        $outageDtoB = $this->createOutage('Outage B');
        $outageA = new Outage(
            $outageDtoA->id,
            new OutagePeriod($outageDtoA->start, $outageDtoA->end),
            new OutageAddress($outageDtoA->streetId, $outageDtoA->streetName, $outageDtoA->buildings, $outageDtoA->city),
            new OutageDescription($outageDtoA->comment)
        );
        $outageB = new Outage(
            $outageDtoB->id,
            new OutagePeriod($outageDtoB->start, $outageDtoB->end),
            new OutageAddress($outageDtoB->streetId, $outageDtoB->streetName, $outageDtoB->buildings, $outageDtoB->city),
            new OutageDescription($outageDtoB->comment)
        );

        // First run
        $this->notifier->handle([$user], [$outageA, $outageB]);

        self::assertCount(1, $this->sender->sent);
        self::assertCount(1, $this->userRepo->saved);
        self::assertEquals('Outage A', $this->sender->sent[0]->comment);

        $this->sender->sent = [];
        $updatedUser = $this->userRepo->saved[0];
        $this->userRepo->saved = [];

        // Second run.
        $this->notifier->handle([$updatedUser], [$outageA, $outageB]);

        self::assertCount(0, $this->sender->sent);
        self::assertCount(0, $this->userRepo->saved);
    }

    private function createOutage(string $comment): OutageDTO
    {
        static $id = 0;
        return new OutageDTO(
            id: ++$id,
            start: new \DateTimeImmutable('2024-11-28T06:47:00+00:00'),
            end: new \DateTimeImmutable('2024-11-28T10:00:00+00:00'),
            city: 'Львів',
            streetId: 12783,
            streetName: 'Шевченка Т.',
            buildings: array_map('trim', explode(',', self::TEST_BUILDINGS)),
            comment: $comment
        );
    }
}
