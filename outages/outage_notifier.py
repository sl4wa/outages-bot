import logging

from telegram import Bot
from telegram.error import Forbidden

from users import UserStorage

from .outage_processor import OutageProcessor


class OutageNotifier:
    def __init__(
        self,
        logger: logging.Logger,
        bot: Bot,
        user_storage: UserStorage,
        outage_processor: OutageProcessor,
    ):
        self.logger = logger
        self.bot = bot
        self.user_storage = user_storage
        self.outage_processor = outage_processor

    async def notify(self):
        for chat_id, user in self.user_storage.all():
            outage = self.outage_processor.get_user_outage(user)

            if outage:
                if user.is_notified(outage):
                    self.logger.info(f"Outage already notified for user {chat_id} - {user.street_name}, {user.building}")
                    continue

                try:
                    await self.bot.send_message(chat_id=chat_id, text=outage.format_message(), parse_mode="HTML")
                    self.user_storage.save(chat_id, user.set_outage(outage))
                    self.logger.info(f"Notification sent to {chat_id} - {user.street_name}, {user.building}")
                except Forbidden:
                    self.user_storage.remove(chat_id)
                    self.logger.info(f"Subscription removed for blocked user {chat_id}.")
                except Exception as e:
                    self.logger.error(f"Failed to send message to {chat_id}: {e}")
            else:
                self.logger.info(f"No relevant outage found for user {chat_id} - {user.street_name}, {user.building}")
