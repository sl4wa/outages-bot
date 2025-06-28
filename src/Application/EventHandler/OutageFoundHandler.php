<?php

namespace App\Application\EventHandler;

use App\Domain\Event\OutageFound;
use App\Domain\Event\OutageNotificationCreated;
use App\Domain\Service\OutageProcessor;
use Symfony\Component\EventDispatcher\Attribute\AsEventListener;
use Symfony\Contracts\EventDispatcher\EventDispatcherInterface;

final class OutageFoundHandler
{
    public function __construct(
        private readonly OutageProcessor $outageProcessor,
        private readonly EventDispatcherInterface $eventDispatcher,
    ) {}

    #[AsEventListener(event: OutageFound::class)]
    public function __invoke(OutageFound $event): void
    {
        $outage = $event->outage;
        $usersToBeChecked = $event->usersToBeChecked;

        $usersToBeNotified = $this->outageProcessor->process($outage, $usersToBeChecked);

        if ($usersToBeNotified) {
            $notificationEvent = new OutageNotificationCreated($outage, $usersToBeNotified);
            $this->eventDispatcher->dispatch($notificationEvent);
        }
    }
}
