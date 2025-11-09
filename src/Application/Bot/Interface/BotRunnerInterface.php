<?php

declare(strict_types=1);

namespace App\Application\Bot\Interface;

interface BotRunnerInterface
{
    /**
     * Start the bot and begin listening for messages/commands
     */
    public function run(): void;
}
