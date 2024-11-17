import os
from dotenv import load_dotenv
from telegram import Bot
from telegram.error import Forbidden
from users import user_storage
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

    async def send_message(self, chat_id: int, message: str) -> None:
        """Send a message to the specified Telegram chat ID."""
        try:
            await self.bot.send_message(chat_id=chat_id, text=message, parse_mode="HTML")
        except Forbidden:
            # Handle case when the bot is blocked by the user
            subscription = user_storage.load_subscription(chat_id)
            if subscription:
                user_storage.remove_subscription(chat_id)
        except Exception as e:
            print(f"Failed to send message to chat_id={chat_id}: {e}")
