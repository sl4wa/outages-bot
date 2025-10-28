<?php
namespace App\Infrastructure\Telegram\Notification;

use App\Application\Notifier\DTO\NotificationSenderDTO;

class TelegramNotificationFormatter
{
    public function format(NotificationSenderDTO $dto): string
    {
        $buildings = implode(', ', $dto->buildings);

        return "Поточні відключення:\n"
            ."Місто: {$dto->city}\n"
            ."Вулиця: {$dto->streetName}\n"
            ."<b>{$dto->start->format('Y-m-d H:i')} – {$dto->end->format('Y-m-d H:i')}</b>\n"
            ."Коментар: {$dto->comment}\n"
            ."Будинки: {$buildings}";
    }
}
