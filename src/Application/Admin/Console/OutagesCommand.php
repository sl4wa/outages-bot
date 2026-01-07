<?php

declare(strict_types=1);

namespace App\Application\Admin\Console;

use App\Application\Interface\OutageProviderInterface;
use Symfony\Component\Console\Attribute\AsCommand;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Helper\Table;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

#[AsCommand(
    name: 'app:outages',
    description: 'Prints a table of outages fetched from the remote API for debug purposes.'
)]
final class OutagesCommand extends Command
{
    public function __construct(
        private readonly OutageProviderInterface $outageProvider
    ) {
        parent::__construct();
    }

    protected function execute(InputInterface $input, OutputInterface $output): int
    {
        $outages = $this->outageProvider->fetchOutages();

        if (!$outages) {
            $output->writeln('<comment>No outages found.</comment>');

            return Command::SUCCESS;
        }

        $table = new Table($output);
        $table->setStyle('compact');
        $table->setColumnMaxWidth(1, 30);
        $table->setColumnMaxWidth(2, 40);
        $table->setHeaders(['StreetID', 'Street', 'Buildings', 'Period', 'Comment']);

        foreach ($outages as $outage) {
            $table->addRow([
                $outage->streetId,
                $outage->streetName,
                implode(', ', $outage->buildings),
                PeriodFormatter::format($outage->start, $outage->end),
                $outage->comment,
            ]);
        }

        $table->render();

        return Command::SUCCESS;
    }
}
