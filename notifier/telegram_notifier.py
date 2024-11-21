import os
from datetime import datetime

from dotenv import load_dotenv
from telegram import Bot
from telegram.error import Forbidden

from users import users

from .notifier_interface import NotifierInterface

load_dotenv()

TELEGRAM_TOKEN_ENV = "TELEGRAM_BOT_TOKEN"


class TelegramNotifier(NotifierInterface):
    """A notifier that sends messages via a Telegram bot."""

    def __init__(self):
        self.token = self._load_bot_token()
        self.bot = Bot(token=self.token)

    def _load_bot_token(self) -> str:
        """Load the Telegram bot token from environment variables."""
        token = os.getenv(TELEGRAM_TOKEN_ENV)
        if not token:
            raise ValueError(
                "No Telegram bot token provided! Set TELEGRAM_BOT_TOKEN in your .env file."
            )
        return token

    def _format_datetime(self, iso_string: str) -> str:
        """Formats the ISO 8601 date string into a readable format."""
        try:
            dt = datetime.fromisoformat(iso_string)
            return dt.strftime("%Y-%m-%d %H:%M")
        except ValueError:
            return iso_string

    async def send_message(self, chat_id: int, relevant_outage) -> None:
        """Send a message to the specified Telegram chat ID."""
        try:
            start_time = self._format_datetime(relevant_outage["dateEvent"])
            end_time = self._format_datetime(relevant_outage["datePlanIn"])
            message = (
                f"Поточні відключення:\n"
                f"Місто: {relevant_outage['city']['name']}\n"
                f"Вулиця: {relevant_outage['street']['name']}\n"
                f"<b>{start_time} - {end_time}</b>\n"
                f"Коментар: {relevant_outage['koment']}\n"
                f"Будинки: {relevant_outage['buildingNames']}"
            )

            await self.bot.send_message(
                chat_id=chat_id, text=message, parse_mode="HTML"
            )
        except Forbidden:
            # Handle case when the bot is blocked by the user
            subscription = users.get(chat_id)
            if subscription:
                users.remove(chat_id)
                logging.info(
                    f"Subscription removed for blocked user {chat_id}.")
