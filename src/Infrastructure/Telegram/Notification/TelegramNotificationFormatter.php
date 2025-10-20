<?php
namespace App\Infrastructure\Telegram\Notification;

use App\Application\DTO\NotificationSenderDTO;

class TelegramNotificationFormatter
{
    public function format(NotificationSenderDTO $dto): string
    {
        $buildings = implode(', ', $dto->outage->address->buildings);

        return "Поточні відключення:\n"
            ."Місто: {$dto->outage->address->city}\n"
            ."Вулиця: {$dto->outage->address->streetName}\n"
            ."<b>{$dto->outage->start->format('Y-m-d H:i')} – {$dto->outage->end->format('Y-m-d H:i')}</b>\n"
            ."Коментар: {$dto->outage->comment}\n"
            ."Будинки: {$buildings}";
    }
}
