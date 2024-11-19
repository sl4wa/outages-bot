# Poweron-Bot

## Setup

### 1. Install Dependencies

Create a virtual environment and install required Python packages:

    make deps

### 2. Install and Enable Services

Set up and enable the bot service using Supervisor:

    make install

This command registers supervisor service:
- `bot`: Manages bot interactions.

### 3. Start the Bot Services

To start service run:

    make start

### 4. Run `notifier.sh` with Cron

To check the API at regular intervals, add `notifier.sh` to cron:

1. Make the script executable:

    ```bash
    chmod +x /path/to/notifier.sh
    ```

2. Add it to cron (for example run every 5 minutes between 7 AM and 11 PM):

    ```bash
    */5 7-23 * * * /path/to/notifier.sh
    ```