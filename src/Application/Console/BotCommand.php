<?php

declare(strict_types=1);

namespace App\Application\Console;

use App\Application\Bot\Interface\BotRunnerInterface;
use Symfony\Component\Console\Attribute\AsCommand;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

#[AsCommand(
    name: 'app:bot',
    description: 'Run the Telegram bot.'
)]
final class BotCommand extends Command
{
    public function __construct(
        private readonly BotRunnerInterface $botRunner
    ) {
        parent::__construct();
    }

    protected function execute(InputInterface $input, OutputInterface $output): int
    {
        $output->writeln('<info>Starting Telegram bot...</info>');
        $this->botRunner->run();

        return Command::SUCCESS;
    }
}
