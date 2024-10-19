# Poweron-Bot

## Setup

### 1. Install dependencies

    make setup

### 2. Run the bot

    ./venv/bin/python3 bot.py

### 3. Add `run_notify.sh` to cron

Make `run_notify.sh` executable:

    chmod +x run_notify.sh

Then add it to cron:

    */5 7-23 * * * /path/to/run_notify.sh
