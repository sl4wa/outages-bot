# Running with Supervisord

This guide explains how to run the outages-bot using supervisord instead of Docker containers.

## Prerequisites

1. **PHP 8.2 or higher** - Check with `php --version`
2. **Composer** - Check with `composer --version`
3. **Supervisor** - Install on your system:
   ```bash
   sudo pacman -S supervisor  # Arch
   sudo apt install supervisor  # Ubuntu/Debian
   ```

## Setup Steps

### 1. Install Dependencies

```bash
composer install
```

### 2. Create Required Directories

```bash
mkdir -p var/log
mkdir -p var/run
```

### 3. Configure Environment

Make sure your `.env`, `.env.local`, and `.env.dev` files are properly configured with all required values (API keys, bot token, etc.).

## Running Supervisord

### Manual Start

```bash
supervisord -c supervisord.conf
```

### As System Service

```bash
sudo systemctl edit supervisor
```

Replace `ExecStart`:

```ini
[Service]
ExecStart=
ExecStart=/usr/bin/supervisord -c /path/to/supervisord.conf
```

Then:

```bash
sudo systemctl daemon-reload
sudo systemctl restart supervisor
sudo systemctl enable supervisor
```

## Management

All management is done via `supervisorctl`. Run commands in the project directory:

```bash
# Check status of all programs
supervisorctl -c supervisord.conf status

# Start a program
supervisorctl -c supervisord.conf start outages-bot
supervisorctl -c supervisord.conf start outages-notifier

# Stop a program
supervisorctl -c supervisord.conf stop outages-bot
supervisorctl -c supervisord.conf stop outages-notifier

# Restart a program
supervisorctl -c supervisord.conf restart outages-bot
supervisorctl -c supervisord.conf restart outages-notifier

# Restart all programs
supervisorctl -c supervisord.conf restart all

# Stop supervisord entirely
supervisorctl -c supervisord.conf shutdown
```

## Viewing Logs

Logs are stored in `var/log/`:

```bash
# Tail the main supervisord log
tail -f var/log/supervisord.log

# Tail the bot log
tail -f var/log/bot.log

# Tail the notifier log
tail -f var/log/notifier.log

# View errors
tail -f var/log/bot.err.log
tail -f var/log/notifier.err.log
```

