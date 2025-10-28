<?php
namespace App\Application\Console;

use App\Application\Notifier\Interface\Provider\OutageProviderInterface;
use Symfony\Component\Console\Attribute\AsCommand;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Helper\Table;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

#[AsCommand(
    name: 'app:outages',
    description: 'Prints a table of outages fetched from the remote API for debug purposes.'
)]
class OutagesCommand extends Command
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
        $table->setHeaders(['City', 'Street', 'StreetID', 'Buildings', 'Period', 'Comment']);

        foreach ($outages as $outage) {
            $buildings = is_array($outage->buildings) ? implode(', ', $outage->buildings) : ($outage->building ?? '');

            $table->addRow([
                $outage->city,
                $outage->streetName ?? $outage->street ?? '',
                $outage->streetId ?? $outage->street ?? '',
                $buildings,
                sprintf(
                    "%s - %s",
                    $outage->start instanceof \DateTimeInterface ? $outage->start->format('Y-m-d H:i') : $outage->start,
                    $outage->end instanceof \DateTimeInterface ? $outage->end->format('Y-m-d H:i') : $outage->end
                ),
                $outage->comment ?? '',
            ]);
        }

        $table->render();

        return Command::SUCCESS;
    }
}
