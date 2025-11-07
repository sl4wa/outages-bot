<?php

declare(strict_types=1);

namespace App\Infrastructure\Dumper;

use App\Application\Interface\DumperInterface;
use Symfony\Component\DependencyInjection\ParameterBag\ParameterBagInterface;
use Symfony\Component\HttpKernel\KernelInterface;
use Symfony\Component\Serializer\SerializerInterface;

readonly class JsonDumper implements DumperInterface
{
    public function __construct(
        private KernelInterface $kernel,
        private SerializerInterface $serializer,
        private ParameterBagInterface $params,
    ) {
    }

    public function dump(mixed $data, string $filename): void
    {
        if (!$this->kernel->isDebug()) {
            return;
        }

        try {
            $debugDir = $this->params->get('kernel.project_dir') . '/data';

            if (!is_dir($debugDir)) {
                mkdir($debugDir, 0770, true);
            }

            $json = $this->serializer->serialize($data, 'json', [
                'json_encode_options' => JSON_PRETTY_PRINT | JSON_UNESCAPED_UNICODE | JSON_UNESCAPED_SLASHES,
            ]);
            file_put_contents($debugDir . '/' . $filename, $json);
        } catch (\Throwable) {
            // Silently fail to prevent breaking the application
        }
    }
}
