<?php

namespace App\Infrastructure\Telegram\Handlers;

use App\Application\Service\UserSubscriptionQueryService;
use SergiX44\Nutgram\Handlers\Type\Command;
use SergiX44\Nutgram\Nutgram;

class SubscriptionInfoCommand extends Command
{
    protected string $command = 'subscription';
    protected ?string $description = 'Показати поточну підписку';

    public function __construct(private readonly UserSubscriptionQueryService $queryService)
    {
        parent::__construct();
    }

    public function handle(Nutgram $bot): void
    {
        $sub = $this->queryService->get($bot->chatId());
        if ($sub) {
            $bot->sendMessage("Ваша поточна підписка:\nВулиця: {$sub->streetName}\nБудинок: {$sub->building}");
        } else {
            $bot->sendMessage("Ви не маєте активної підписки.");
        }
    }
}
