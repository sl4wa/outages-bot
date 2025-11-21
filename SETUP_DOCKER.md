# Running with Docker

This guide explains how to run the outages-bot using Docker containers.

## Prerequisites

1. **Docker** - Install from https://docs.docker.com/install/
2. **Docker Compose** - Install from https://docs.docker.com/compose/install/

## Setup Steps

### 1. Configure Environment

Make sure your `.env`, `.env.local`, and `.env.dev` files are properly configured with all required values (API keys, bot token, etc.).

## Running Docker

Start the containers:

```bash
docker-compose up -d
```

This starts both services in the background:
- `outages-bot`: The main Telegram bot
- `outages-notifier`: Cron job that runs every 5 minutes

## Management

```bash
# Check status of containers
docker-compose ps

# View logs from all services
docker-compose logs -f

# View logs from specific service
docker-compose logs -f bot
docker-compose logs -f notifier

# Stop all services
docker-compose stop

# Start all services
docker-compose start

# Restart all services
docker-compose restart

# Stop and remove containers
docker-compose down
```

## Rebuilding After Changes

If you modify dependencies:

```bash
docker-compose build
docker-compose up -d
```
