<?php

declare(strict_types=1);

namespace App\Application\Admin\Console;

use App\Application\Interface\TelegramUserInfoProviderInterface;
use App\Domain\Repository\UserRepositoryInterface;
use Symfony\Component\Console\Attribute\AsCommand;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Helper\Table;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;
use Throwable;

#[AsCommand(
    name: 'app:users',
    description: 'List all subscribed users with their Telegram info and addresses.'
)]
final class UsersCommand extends Command
{
    public function __construct(
        private readonly UserRepositoryInterface $userRepository,
        private readonly TelegramUserInfoProviderInterface $userInfoProvider,
    ) {
        parent::__construct();
    }

    protected function execute(InputInterface $input, OutputInterface $output): int
    {
        $users = $this->userRepository->findAll();

        if (!$users) {
            $output->writeln('<comment>No users found.</comment>');

            return Command::SUCCESS;
        }

        usort($users, static function ($a, $b) {
            if ($a->outageInfo === null && $b->outageInfo === null) {
                return 0;
            }
            if ($a->outageInfo === null) {
                return 1;
            }
            if ($b->outageInfo === null) {
                return -1;
            }

            return $b->outageInfo->period->startDate <=> $a->outageInfo->period->startDate;
        });

        $table = new Table($output);
        $table->setStyle('compact');
        $table->setHeaders(['Chat ID', 'Username', 'First Name', 'Last Name', 'Street', 'Building', 'Outage', 'Comment']);

        $successCount = 0;

        foreach ($users as $user) {
            try {
                $userInfo = $this->userInfoProvider->getUserInfo($user->id);

                $table->addRow([
                    $userInfo->chatId,
                    $userInfo->username ? '@' . $userInfo->username : '-',
                    $this->sanitize($userInfo->firstName),
                    $this->sanitize($userInfo->lastName),
                    $user->address->streetName,
                    $user->address->building,
                    $user->outageInfo !== null
                        ? PeriodFormatter::format($user->outageInfo->period->startDate, $user->outageInfo->period->endDate)
                        : '-',
                    $user->outageInfo?->description->value ?? '-',
                ]);

                ++$successCount;
            } catch (Throwable $e) {
                $output->writeln("<error>Failed to get info for chat {$user->id}: {$e->getMessage()}</error>");
            }
        }

        $table->render();

        $output->writeln('');
        $output->writeln("<info>Total Users: {$successCount}</info>");

        return Command::SUCCESS;
    }

    private function sanitize(?string $value): string
    {
        if (!$value) {
            return '-';
        }

        // Remove invisible/control characters
        $cleaned = preg_replace('/[\p{Cf}\x{3164}]/u', '', $value);
        $trimmed = trim((string) $cleaned);

        return $trimmed !== '' ? $trimmed : '-';
    }
}
