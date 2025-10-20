<?php
namespace App\Application\Command;

use App\Application\Interface\Repository\UserRepositoryInterface;
use App\Application\Service\NotifierService;
use App\Application\Service\OutageFetchService;
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
    public function __construct(
        private readonly NotifierService $notificationService,
        private readonly OutageFetchService $outageFetchService,
        private readonly UserRepositoryInterface $userRepository,
    ) {
        parent::__construct();
    }

    protected function execute(InputInterface $input, OutputInterface $output): int
    {
        $outages = $this->outageFetchService->fetch();
        $users = $this->userRepository->findAll();

        $sent = $this->notificationService->notify($users, $outages);
        $output->writeln("<info>Successfully dispatched $sent outages.</info>");
        return Command::SUCCESS;
    }
}
