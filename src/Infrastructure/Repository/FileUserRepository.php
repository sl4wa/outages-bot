<?php

namespace App\Infrastructure\Repository;

use App\Application\Interface\Repository\UserRepositoryInterface;
use App\Domain\Entity\User;
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
            'street_id'   => $user->streetId,
            'street_name' => $user->streetName,
            'building'    => $user->building,
            'start_date'  => $user->startDate instanceof \DateTimeImmutable ? $user->startDate->format(DATE_ATOM) : '',
            'end_date'    => $user->endDate instanceof \DateTimeImmutable ? $user->endDate->format(DATE_ATOM) : '',
            'comment'     => $user->comment,
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
        return new User(
            $id,
            (int)$fields['street_id'],
            $fields['street_name'],
            $fields['building'],
            $fields['start_date'] ? new \DateTimeImmutable($fields['start_date']) : null,
            $fields['end_date'] ? new \DateTimeImmutable($fields['end_date']) : null,
            $fields['comment']
        );
    }
}
