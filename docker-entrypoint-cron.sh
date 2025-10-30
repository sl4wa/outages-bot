#!/bin/sh
set -eu

CRON_FILE=/etc/crontabs/root
mkdir -p /etc/crontabs

: "${APP_ENV:=prod}"

{
  echo "SHELL=/bin/sh"
  echo "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
  echo "APP_ENV=${APP_ENV}"
  if [ -n "${TELEGRAM_BOT_TOKEN:-}" ]; then
    echo "TELEGRAM_BOT_TOKEN=${TELEGRAM_BOT_TOKEN}"
  fi
  echo "*/10 * * * * cd /app && php bin/console app:notifier >> /dev/stdout 2>&1"
} > "$CRON_FILE"

exec crond -f -l 8
