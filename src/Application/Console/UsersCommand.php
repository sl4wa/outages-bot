<?php

declare(strict_types=1);

namespace App\Application\Console;

use App\Application\Interface\Repository\UserRepositoryInterface;
use SergiX44\Nutgram\Nutgram;
use Symfony\Component\Console\Attribute\AsCommand;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Helper\Table;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

#[AsCommand(
    name: 'app:users',
    description: 'List all subscribed users with their Telegram info and addresses.'
)]
final class UsersCommand extends Command
{
    public function __construct(
        private readonly UserRepositoryInterface $userRepository,
        private readonly Nutgram $nutgram,
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

        $table = new Table($output);
        $table->setHeaders(['Chat ID', 'Username', 'First Name', 'Last Name', 'Street', 'Building']);

        $successCount = 0;

        foreach ($users as $user) {
            try {
                $chat = $this->nutgram->getChat($user->id);

                $table->addRow([
                    $chat->id,
                    $chat->username ? '@' . $chat->username : '-',
                    $this->sanitize($chat->first_name),
                    $this->sanitize($chat->last_name),
                    $user->address->streetName,
                    $user->address->building,
                ]);

                $successCount++;
            } catch (\Throwable $e) {
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
        $trimmed = trim($cleaned);

        return $trimmed !== '' ? $trimmed : '-';
    }
}
