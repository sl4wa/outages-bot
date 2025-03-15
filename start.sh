#!/bin/sh

# Export the timezone so cron (and its children) see it.
export TZ=Europe/Kyiv

# Start cron in the background.
nohup /usr/sbin/cron -f &

# Tail the notifier log file so cron output is sent to Docker's stdout.
tail -F /app/cron.log &

# Start the bot process in the foreground.
/usr/local/bin/python /app/bot.py
