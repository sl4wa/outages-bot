import os
from datetime import datetime
from dotenv import load_dotenv
from telegram import Bot
from telegram.error import Forbidden
from .outage import Outage

TELEGRAM_TOKEN_ENV = "TELEGRAM_BOT_TOKEN"


class OutagesNotifier:
    """A notifier that sends messages via a Telegram bot."""

    def __init__(self):
        self.token = self._load_bot_token()
        self.bot = Bot(token=self.token)

    def _load_bot_token(self) -> str:
        """Load the Telegram bot token from environment variables."""
        load_dotenv()
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

    async def send_message(self, chat_id: int, outage: Outage) -> None:
        """Send a message to the specified Telegram chat ID."""
        start = self._format_datetime(outage.start_date)
        end = self._format_datetime(outage.end_date)
        message = (
            f"Поточні відключення:\n"
            f"Місто: {outage.city}\n"
            f"Вулиця: {outage.street}\n"
            f"<b>{start} - {end}</b>\n"
            f"Коментар: {outage.comment}\n"
            f"Будинки: {outage.building}"
        )

        await self.bot.send_message(chat_id=chat_id, text=message, parse_mode="HTML")
