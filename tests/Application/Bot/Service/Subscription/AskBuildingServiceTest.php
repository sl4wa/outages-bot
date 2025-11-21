<?php

declare(strict_types=1);

namespace App\Tests\Application\Bot\Service\Subscription;

use App\Application\Bot\Command\CreateOrUpdateUserSubscriptionCommandHandler;
use App\Application\Bot\Service\Subscription\AskBuildingService;
use App\Application\Interface\Repository\UserRepositoryInterface;
use PHPUnit\Framework\TestCase;

final class AskBuildingServiceTest extends TestCase
{
    private AskBuildingService $service;
    private UserRepositoryInterface $userRepository;

    protected function setUp(): void
    {
        $this->userRepository = $this->createMock(UserRepositoryInterface::class);
        $commandHandler = new CreateOrUpdateUserSubscriptionCommandHandler($this->userRepository);
        $this->service = new AskBuildingService($commandHandler);
    }

    public function testSuccessfulBuildingSubmission(): void
    {
        $this->userRepository
            ->expects($this->once())
            ->method('save');

        $result = $this->service->handle(
            chatId: 12345,
            selectedStreetId: 123,
            selectedStreetName: 'Шевченка',
            building: '13'
        );

        self::assertTrue($result->isSuccess);
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
            selectedStreetId: 123,
            selectedStreetName: 'Шевченка',
            building: '196-А'
        );

        self::assertTrue($result->isSuccess);
        self::assertStringContainsString('196-А', $result->message);
    }

    public function testEmptyBuildingError(): void
    {
        $this->userRepository
            ->expects($this->never())
            ->method('save');

        $result = $this->service->handle(
            chatId: 12345,
            selectedStreetId: 123,
            selectedStreetName: 'Шевченка',
            building: ''
        );

        self::assertFalse($result->isSuccess);
        self::assertSame('Введіть номер будинку.', $result->message);
    }

    public function testWhitespaceOnlyBuildingError(): void
    {
        $this->userRepository
            ->expects($this->never())
            ->method('save');

        $result = $this->service->handle(
            chatId: 12345,
            selectedStreetId: 123,
            selectedStreetName: 'Шевченка',
            building: '   '
        );

        self::assertFalse($result->isSuccess);
        self::assertSame('Введіть номер будинку.', $result->message);
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
            selectedStreetId: 123,
            selectedStreetName: 'Шевченка',
            building: $building
        );

        self::assertFalse($result->isSuccess);
        self::assertSame('Невірний формат номера будинку', $result->message);
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

    public function testMissingStreetIdError(): void
    {
        $this->userRepository
            ->expects($this->never())
            ->method('save');

        $result = $this->service->handle(
            chatId: 12345,
            selectedStreetId: null,
            selectedStreetName: 'Шевченка',
            building: '13'
        );

        self::assertFalse($result->isSuccess);
        self::assertSame('Підписка не завершена. Будь ласка, почніть знову.', $result->message);
    }

    public function testMissingStreetNameError(): void
    {
        $this->userRepository
            ->expects($this->never())
            ->method('save');

        $result = $this->service->handle(
            chatId: 12345,
            selectedStreetId: 123,
            selectedStreetName: null,
            building: '13'
        );

        self::assertFalse($result->isSuccess);
        self::assertSame('Підписка не завершена. Будь ласка, почніть знову.', $result->message);
    }

    public function testMissingBothStreetIdAndNameError(): void
    {
        $this->userRepository
            ->expects($this->never())
            ->method('save');

        $result = $this->service->handle(
            chatId: 12345,
            selectedStreetId: null,
            selectedStreetName: null,
            building: '13'
        );

        self::assertFalse($result->isSuccess);
        self::assertSame('Підписка не завершена. Будь ласка, почніть знову.', $result->message);
    }
}
