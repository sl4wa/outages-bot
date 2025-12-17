<?php

declare(strict_types=1);

namespace App\Application\Notifier\Console;

use App\Application\Interface\DumperInterface;
use App\Application\Notifier\Service\NotificationService;
use App\Application\Notifier\Service\OutageFetchService;
use Symfony\Component\Console\Attribute\AsCommand;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

#[AsCommand(
    name: 'app:notifier',
    description: 'Send outage notifications to users.',
)]
final class NotifierCommand extends Command
{
    public function __construct(
        private readonly NotificationService $notificationService,
        private readonly OutageFetchService $outageFetchService,
        private readonly DumperInterface $dumper,
    ) {
        parent::__construct();
    }

    protected function execute(InputInterface $input, OutputInterface $output): int
    {
        $outages = $this->outageFetchService->handle();
        $this->dumper->dump($outages, 'outages.json');
        $sent = $this->notificationService->handle($outages);
        $output->writeln("<info>Successfully dispatched $sent outages.</info>");

        return Command::SUCCESS;
    }
}
