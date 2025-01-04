import logging

from telegram import Bot
from telegram.error import Forbidden

from env import load_bot_token
from users import UserStorage

from .outage_processor import OutageProcessor


class OutageNotifier:
    def __init__(self, logger: logging.Logger):
        self.logger = logger

    async def notify(self):
        bot = Bot(token=load_bot_token())

        outage_processor = OutageProcessor()
        user_storage = UserStorage()

        for chat_id, user in user_storage.all():
            outage = outage_processor.get_user_outage(user)

            if outage:
                if user.is_notified(outage):
                    self.logger.info(f"Outage already notified for user {chat_id} - {user.street_name}, {user.building}")
                    continue

                try:
                    await bot.send_message(chat_id=chat_id, text=outage.format_message(), parse_mode="HTML")
                    user_storage.save(chat_id, user.set_outage(outage))
                    self.logger.info(f"Notification sent to {chat_id} - {user.street_name}, {user.building}")
                except Forbidden:
                    user_storage.remove(chat_id)
                    self.logger.info(f"Subscription removed for blocked user {chat_id}.")
                except Exception as e:
                    self.logger.error(f"Failed to send message to {chat_id}: {e}")
            else:
                self.logger.info(f"No relevant outage found for user {chat_id} - {user.street_name}, {user.building}")
