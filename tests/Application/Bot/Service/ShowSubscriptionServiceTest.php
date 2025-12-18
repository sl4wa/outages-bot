<?php

declare(strict_types=1);

namespace App\Tests\Application\Bot\Service;

use App\Application\Bot\Query\GetUserSubscriptionQueryHandler;
use App\Application\Bot\Service\ShowSubscriptionService;
use App\Domain\Entity\User;
use App\Domain\Repository\UserRepositoryInterface;
use App\Domain\ValueObject\UserAddress;
use PHPUnit\Framework\TestCase;
use RuntimeException;

final class ShowSubscriptionServiceTest extends TestCase
{
    private ShowSubscriptionService $service;

    private UserRepositoryInterface $userRepository;

    protected function setUp(): void
    {
        $this->userRepository = $this->createMock(UserRepositoryInterface::class);
        $queryHandler = new GetUserSubscriptionQueryHandler($this->userRepository);
        $this->service = new ShowSubscriptionService($queryHandler);
    }

    public function testReturnsSimplePromptForNewUser(): void
    {
        $this->userRepository
            ->expects($this->once())
            ->method('find')
            ->with(12345)
            ->willReturn(null);

        $result = $this->service->handle(12345);

        self::assertSame('Будь ласка, введіть назву вулиці:', $result);
    }

    public function testReturnsPromptWithCurrentSubscriptionForExistingUser(): void
    {
        $user = new User(
            id: 12345,
            address: new UserAddress(
                streetId: 12783,
                streetName: 'Шевченка Т.',
                building: '271'
            ),
            outageInfo: null
        );

        $this->userRepository
            ->expects($this->once())
            ->method('find')
            ->with(12345)
            ->willReturn($user);

        $result = $this->service->handle(12345);

        self::assertStringContainsString('Ваша поточна підписка:', $result);
        self::assertStringContainsString('Шевченка Т.', $result);
        self::assertStringContainsString('271', $result);
        self::assertStringContainsString('нову назву вулиці для оновлення підписки', $result);
    }

    public function testHandlesDifferentChatIds(): void
    {
        $this->userRepository
            ->expects($this->once())
            ->method('find')
            ->with(99999)
            ->willReturn(null);

        $result = $this->service->handle(99999);

        self::assertStringContainsString('Будь ласка, введіть назву вулиці:', $result);
    }

    public function testShowsCorrectStreetNameInSubscription(): void
    {
        $user = new User(
            id: 12345,
            address: new UserAddress(
                streetId: 12444,
                streetName: 'Молдавська',
                building: '13-А'
            ),
            outageInfo: null
        );

        $this->userRepository
            ->expects($this->once())
            ->method('find')
            ->with(12345)
            ->willReturn($user);

        $result = $this->service->handle(12345);

        self::assertStringContainsString('Молдавська', $result);
        self::assertStringContainsString('13-А', $result);
    }

    public function testMessageFormatWithCyrillicCharacters(): void
    {
        $user = new User(
            id: 12345,
            address: new UserAddress(
                streetId: 12783,
                streetName: 'Київська',
                building: '196-А'
            ),
            outageInfo: null
        );

        $this->userRepository
            ->expects($this->once())
            ->method('find')
            ->with(12345)
            ->willReturn($user);

        $result = $this->service->handle(12345);

        self::assertStringContainsString('Вулиця: Київська', $result);
        self::assertStringContainsString('Будинок: 196-А', $result);
    }

    public function testHandlesCorruptedUserDataGracefully(): void
    {
        $this->userRepository
            ->expects($this->once())
            ->method('find')
            ->with(12345)
            ->willThrowException(new RuntimeException('Corrupted data'));

        $result = $this->service->handle(12345);

        self::assertSame('Будь ласка, введіть назву вулиці:', $result);
    }
}
