<?php
namespace App\Domain\Entity;

class Outage
{
    public function __construct(
        public readonly \DateTimeImmutable $start,
        public readonly \DateTimeImmutable $end,
        public readonly string $city,
        public readonly int $streetId,
        public readonly string $streetName,
        public readonly array $buildingNames,
        public readonly string $comment
    ) {}

    public static function fromArray(array $data): self
    {
        $comment = $data['koment'] ?? '';
        $comment = preg_replace('/[\r\n]+/', ' ', $comment);
        $comment = trim($comment);

        return new self(
            new \DateTimeImmutable($data['dateEvent']),
            new \DateTimeImmutable($data['datePlanIn']),
            $data['city']['name'],
            (int)$data['street']['id'],
            $data['street']['name'],
            array_map('trim', explode(',', $data['buildingNames'])),
            $comment
        );
    }

    public function matchesUser(User $user): bool
    {
        return $this->streetId === $user->streetId &&
            in_array($user->building, $this->buildingNames, true);
    }

    public function isIdenticalPeriodAndComment(User $user): bool
    {
        return $user->startDate == $this->start
            && $user->endDate == $this->end
            && $user->comment == $this->comment;
    }
}
