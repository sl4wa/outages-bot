# Poweron Outages Notifier

Telegram bot for power outages in Lviv, Ukraine.

## Console Commands

- `bin/console app:bot`  
  Runs the Telegram bot for managing subscriptions.

- `bin/console app:notifier`  
  Cron command: check if any user has relevant outage.

- `bin/console app:outages`  
  List current outages from API.

- `bin/console app:users`
  List all subscribed users with their Telegram info and addresses.

## TODO

- [x] tests
- [ ] rector
- [ ] csfix
- [ ] phpstan
- [x] docker
- [ ] logging
