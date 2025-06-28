<?php

namespace App\Infrastructure\Telegram\Bot;

use Symfony\Component\Console\Attribute\AsCommand;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

#[AsCommand(
    name: 'app:bot',
    description: 'Run the Telegram bot using Nutgram.'
)]
class BotCommand extends Command
{
    private TelegramBotRunner $telegramBotRunner;

    public function __construct(TelegramBotRunner $telegramBotRunner)
    {
        parent::__construct();
        $this->telegramBotRunner = $telegramBotRunner;
    }

    protected function execute(InputInterface $input, OutputInterface $output): int
    {
        $output->writeln('<info>Starting Telegram bot (Nutgram)...</info>');
        $this->telegramBotRunner->run();
        return Command::SUCCESS;
    }
}
