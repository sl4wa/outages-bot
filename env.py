import logging
import os

from dotenv import load_dotenv


def load_bot_token() -> str:
    """Load the Telegram bot token from environment variables."""
    load_dotenv()
    token = os.getenv("TELEGRAM_BOT_TOKEN")
    
    if not token:
        error_message = (
            "No token provided! Please add the TELEGRAM_BOT_TOKEN to the .env file as follows:\n\n"
            "TELEGRAM_BOT_TOKEN=your-telegram-bot-token"
        )
        logging.error(error_message)
        raise ValueError(error_message)
    
    return token
