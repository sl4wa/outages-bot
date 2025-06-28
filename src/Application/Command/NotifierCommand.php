<?php
namespace App\Application\Command;

use App\Application\Service\NotifierService;
use Symfony\Component\Console\Attribute\AsCommand;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

#[AsCommand(
    name: 'app:notifier',
    description: 'Send outage notifications to users.',
)]
class NotifierCommand extends Command
{
    public function __construct(private readonly NotifierService $notificationService)
    {
        parent::__construct();
    }

    protected function execute(InputInterface $input, OutputInterface $output): int
    {
        $sent = $this->notificationService->notify();
        $output->writeln("<info>Successfully dispatched $sent outages.</info>");
        return Command::SUCCESS;
    }
}
