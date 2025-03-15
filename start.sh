#!/bin/sh

# Start cron in the background.
nohup /usr/sbin/cron -f &

# Start the bot process in the foreground.
/usr/local/bin/python /app/bot.py
