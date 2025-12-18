<?php

declare(strict_types=1);

namespace App\Tests\Application\Bot\Service;

use App\Application\Bot\Command\CreateOrUpdateUserSubscriptionCommandHandler;
use App\Application\Bot\Service\SaveSubscriptionService;
use App\Domain\Repository\UserRepositoryInterface;
use PHPUnit\Framework\TestCase;

final class SaveSubscriptionServiceTest extends TestCase
{
    private SaveSubscriptionService $service;

    private UserRepositoryInterface $userRepository;

    protected function setUp(): void
    {
        $this->userRepository = $this->createMock(UserRepositoryInterface::class);
        $commandHandler = new CreateOrUpdateUserSubscriptionCommandHandler($this->userRepository);
        $this->service = new SaveSubscriptionService($commandHandler);
    }

    public function testSuccessfulBuildingSubmission(): void
    {
        $this->userRepository
            ->expects($this->once())
            ->method('save');

        $result = $this->service->handle(
            chatId: 12345,
            streetId: 123,
            streetName: 'Шевченка',
            building: '13'
        );

        self::assertTrue($result->success);
        self::assertStringContainsString('Ви підписалися', $result->message);
        self::assertStringContainsString('Шевченка', $result->message);
        self::assertStringContainsString('13', $result->message);
    }

    public function testSuccessfulBuildingSubmissionWithLetter(): void
    {
        $this->userRepository
            ->expects($this->once())
            ->method('save');

        $result = $this->service->handle(
            chatId: 12345,
            streetId: 123,
            streetName: 'Шевченка',
            building: '196-А'
        );

        self::assertTrue($result->success);
        self::assertStringContainsString('196-А', $result->message);
    }

    public function testEmptyBuildingError(): void
    {
        $this->userRepository
            ->expects($this->never())
            ->method('save');

        $result = $this->service->handle(
            chatId: 12345,
            streetId: 123,
            streetName: 'Шевченка',
            building: ''
        );

        self::assertFalse($result->success);
        self::assertStringContainsString('Невірний формат номера будинку', $result->message);
    }

    public function testWhitespaceOnlyBuildingError(): void
    {
        $this->userRepository
            ->expects($this->never())
            ->method('save');

        $result = $this->service->handle(
            chatId: 12345,
            streetId: 123,
            streetName: 'Шевченка',
            building: '   '
        );

        self::assertFalse($result->success);
        self::assertStringContainsString('Невірний формат номера будинку', $result->message);
    }

    /**
     * @dataProvider invalidBuildingFormatsProvider
     */
    public function testInvalidBuildingFormatError(string $building): void
    {
        $this->userRepository
            ->expects($this->never())
            ->method('save');

        $result = $this->service->handle(
            chatId: 12345,
            streetId: 123,
            streetName: 'Шевченка',
            building: $building
        );

        self::assertFalse($result->success);
        self::assertStringContainsString('Невірний формат номера будинку', $result->message);
    }

    public static function invalidBuildingFormatsProvider(): array
    {
        return [
            'with parentheses' => ['291(0083)'],
            'with comma' => ['350,А'],
            'multiple letters' => ['13-AB'],
            'slash separator' => ['13/A'],
            'only letters' => ['abc'],
            'number after hyphen' => ['13-1'],
        ];
    }
}
