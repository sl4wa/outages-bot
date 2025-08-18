# Poweron Outages Notifier

Telegram bot for power outages in Lviv, Ukraine.

## Console Commands

- `symfony console app:bot`  
  Runs the Telegram bot for managing subscriptions.

- `symfony console app:notifier`  
  Cron command: check if any user has relevant outage.

- `symfony console app:outages`  
  List current outages from API.

- `symfony console app:users`  
  *TODO* List subscribed users

## TODO

- [x] tests
- [ ] rector
- [ ] csfix
- [ ] phpstan
- [x] docker
- [ ] logging
