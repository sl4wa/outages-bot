<?php

namespace App\Infrastructure\Telegram\Handlers;

use App\Application\Bot\Service\GetSubscriptionInfoService;
use SergiX44\Nutgram\Handlers\Type\Command;
use SergiX44\Nutgram\Nutgram;

class SubscriptionInfoCommand extends Command
{
    protected string $command = 'subscription';
    protected ?string $description = 'Показати поточну підписку';

    public function __construct(private readonly GetSubscriptionInfoService $getSubscriptionInfoService)
    {
        parent::__construct();
    }

    public function handle(Nutgram $bot): void
    {
        $bot->sendMessage($this->getSubscriptionInfoService->handle($bot->chatId()));
    }
}
