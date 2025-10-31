<?php

namespace App\Infrastructure\Repository;

use App\Application\Interface\Repository\UserRepositoryInterface;
use App\Domain\Entity\User;
use App\Domain\ValueObject\OutageDescription;
use App\Domain\ValueObject\OutageInfo;
use App\Domain\ValueObject\OutagePeriod;
use App\Domain\ValueObject\UserAddress;
use Symfony\Component\DependencyInjection\ParameterBag\ParameterBagInterface;

class FileUserRepository implements UserRepositoryInterface
{
    private string $dataDir;

    public function __construct(ParameterBagInterface $params)
    {
        $projectDir = $params->get('kernel.project_dir');
        $this->dataDir = $projectDir . '/data/users';

        if (!is_dir($this->dataDir)) {
            mkdir($this->dataDir, 0770, true);
        }
    }

    public function findAll(): array
    {
        $users = [];
        foreach (glob($this->dataDir . '/*.txt') as $file) {
            if ($user = $this->loadFromFile($file)) {
                $users[] = $user;
            }
        }
        return $users;
    }

    public function find(int $chatId): ?User
    {
        $file = $this->getFilePath($chatId);
        return file_exists($file) ? $this->loadFromFile($file) : null;
    }

    public function save(User $user): void
    {
        $fields = [
            'street_id'   => $user->address->streetId,
            'street_name' => $user->address->streetName,
            'building'    => $user->address->building,
            'start_date'  => $user->outageInfo ? $user->outageInfo->period->startDate->format(DATE_ATOM) : '',
            'end_date'    => $user->outageInfo ? $user->outageInfo->period->endDate->format(DATE_ATOM) : '',
            'comment'     => $user->outageInfo ? $user->outageInfo->description->value : '',
        ];
        $lines = [];
        foreach ($fields as $key => $val) {
            $lines[] = "$key: $val";
        }
        file_put_contents($this->getFilePath($user->id), implode(PHP_EOL, $lines));
    }

    public function remove(int $chatId): void
    {
        $file = $this->getFilePath($chatId);
        if (file_exists($file)) {
            unlink($file);
        }
    }

    private function getFilePath(int $chatId): string
    {
        return $this->dataDir . '/' . $chatId . '.txt';
    }

    private function loadFromFile(string $file): ?User
    {
        $id = (int)basename($file, '.txt');
        $fields = [
            'street_id'   => 0,
            'street_name' => '',
            'building'    => '',
            'start_date'  => '',
            'end_date'    => '',
            'comment'     => ''
        ];
        $data = file($file, FILE_IGNORE_NEW_LINES | FILE_SKIP_EMPTY_LINES);
        foreach ($data as $line) {
            if (strpos($line, ':') !== false) {
                [$key, $val] = array_map('trim', explode(':', $line, 2));
                if (array_key_exists($key, $fields)) {
                    $fields[$key] = $val;
                }
            }
        }

        $address = new UserAddress(
            (int)$fields['street_id'],
            $fields['street_name'],
            $fields['building']
        );

        $outageInfo = null;
        if ($fields['start_date'] && $fields['end_date']) {
            $period = new OutagePeriod(
                new \DateTimeImmutable($fields['start_date']),
                new \DateTimeImmutable($fields['end_date'])
            );
            $description = new OutageDescription($fields['comment']);
            $outageInfo = new OutageInfo($period, $description);
        }

        return new User(
            $id,
            $address,
            $outageInfo
        );
    }
}
