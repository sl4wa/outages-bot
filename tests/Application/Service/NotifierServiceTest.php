<?php

declare(strict_types=1);

namespace App\Tests\Application\Service;

use App\Application\Service\NotifierService;
use App\Domain\Entity\Outage;
use App\Domain\Entity\User;
use App\Tests\Support\TestNotificationSender;
use App\Tests\Support\TestOutageProvider;
use App\Tests\Support\TestUserRepository;
use Symfony\Bundle\FrameworkBundle\Test\KernelTestCase;

final class NotifierServiceTest extends KernelTestCase
{
    protected static function getKernelClass(): string
    {
        return \App\Kernel::class;
    }

    private NotifierService $notifier;
    private TestNotificationSender $sender;
    private TestUserRepository $userRepo;
    private TestOutageProvider $provider;

    protected function setUp(): void
    {
        self::bootKernel();

        $container = self::getContainer();

        // Fetch services from the test container (DI)
        $this->sender = $container->get(TestNotificationSender::class);
        $this->userRepo = $container->get(TestUserRepository::class);
        $this->provider = $container->get(TestOutageProvider::class);
        $this->notifier = $container->get(NotifierService::class);

        // Reset state between tests (shared services)
        $this->sender->sent = [];
        $this->sender->blockUserId = null;
        $this->userRepo->all = [];
        $this->userRepo->saved = [];
        $this->userRepo->removed = [];
        $this->provider->outages = [];
    }

    public function testNotificationSentAndUserSaved(): void
    {
        $outage = $this->createOutage('Застосування ГПВ');
        $user = new User(100, 12783, 'Шевченка Т.', '271', null, null, '');

        // prepare doubles
        $this->userRepo->all = [$user];
        $this->provider->outages = [$outage];

        $this->notifier->notify();

        self::assertCount(1, $this->sender->sent); // one notification
        self::assertEquals(100, $this->sender->sent[0]->user->id);
        self::assertCount(1, $this->userRepo->saved);
        self::assertEquals('Застосування ГПВ', $this->userRepo->saved[0]->comment);
    }

    public function testSubscriptionRemovedForBlockedUser(): void
    {
        $this->sender->blockUserId = 101; // simulate Forbidden

        $outage = $this->createOutage('Застосування ГПВ');
        $user = new User(101, 12783, 'Шевченка Т.', '271', null, null, '');

        $this->userRepo->all = [$user];
        $this->provider->outages = [$outage];

        $this->notifier->notify();

        self::assertSame([101], $this->userRepo->removed);
        self::assertCount(0, $this->sender->sent);
        self::assertCount(0, $this->userRepo->saved);
    }

    public function testNoRelevantOutageProducesNoNotification(): void
    {
        $outage = $this->createOutage('Застосування ГПВ');
        $user = new User(102, 99999, 'Nonexistent Street', '1', null, null, '');

        $this->userRepo->all = [$user];
        $this->provider->outages = [$outage];

        $this->notifier->notify();

        self::assertCount(0, $this->sender->sent);
        self::assertCount(0, $this->userRepo->saved);
        self::assertCount(0, $this->userRepo->removed);
    }

    public function testMultipleOutagesSameBuildingNotifiesOnceWithFirstComment(): void
    {
        $user = new User(103, 12783, 'Шевченка Т.', '271', null, null, '');

        $first = $this->createOutage('Застосування ГАВ');
        $second = $this->createOutage('Застосування ГПВ');

        $this->userRepo->all = [$user];
        $this->provider->outages = [$first, $second];

        $this->notifier->notify();

        self::assertCount(1, $this->sender->sent, 'should send only once');
        self::assertCount(1, $this->userRepo->saved, 'should save only once');
        self::assertEquals('Застосування ГАВ', $this->userRepo->saved[0]->comment, 'first outage wins');
    }

    private function createOutage(string $comment): Outage
    {
        $buildings = '271, 273, 273-А, 275, 277, 279, 281, 281-А, 282, 283, 283-А, '
            . '284, 284-А, 285, 285-А, 287, 289, 289-А, 290-А, 291, 291(0083), '
            . '293, 295, 297, 297-А, 297-Б, 308, 313, 316, 316-А, 318, 318-А, '
            . '320, 322, 324, 326, 328, 328-А, 330, 332, 334, 336, 338, 340-А, '
            . '342, 346, 348-А, 350, 350,А, 350-В, 358, 358-А, 360-В';

        return new Outage(
            new \DateTimeImmutable('2024-11-28T06:47:00+00:00'),
            new \DateTimeImmutable('2024-11-28T10:00:00+00:00'),
            'Львів',
            12783,
            'Шевченка Т.',
            array_map('trim', explode(',', $buildings)),
            $comment
        );
    }
}
