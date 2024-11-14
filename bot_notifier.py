import asyncio
import logging
import os
import sys
from logging.handlers import WatchedFileHandler

from dotenv import load_dotenv
from telegram import Bot
from telegram.error import Forbidden

from users import user_storage

# Constants
PIPE_NAME = "telegram.pipe"
LOG_FILE = "bot_notifier.log"
TELEGRAM_TOKEN_ENV = "TELEGRAM_BOT_TOKEN"


def configure_logging() -> None:
    """Configure logging for the application."""
    logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
        datefmt="%Y-%m-%d %H:%M:%S",
        handlers=[WatchedFileHandler(LOG_FILE), logging.StreamHandler(sys.stdout)],
    )

    # Suppress excessive logging from HTTP requests library
    httpx_logger = logging.getLogger("httpx")
    httpx_logger.setLevel(logging.WARNING)

    logging.info("Logging is configured.")


def load_bot_token() -> str:
    """Load the Telegram bot token from environment variables."""
    token = os.getenv(TELEGRAM_TOKEN_ENV)
    if not token:
        error_message = (
            "No token provided! Please add the TELEGRAM_BOT_TOKEN to the .env file as follows:\n\n"
            "TELEGRAM_BOT_TOKEN=your-telegram-bot-token"
        )
        logging.error(error_message)
        raise ValueError(error_message)
    logging.info("Telegram bot token loaded.")
    return token


async def send_message(bot, chat_id: int, text: str) -> None:
    """Send a message to the specified Telegram chat ID."""
    try:
        await bot.send_message(chat_id=chat_id, text=text, parse_mode="HTML")
        logging.info(f"Sent message to chat_id={chat_id}")
    except Forbidden as e:
        logging.warning(f"Bot was blocked by the user with chat_id={chat_id}")
        subscriptions = user_storage.load_subscriptions()
        if chat_id in subscriptions:
            del subscriptions[chat_id]
            user_storage.save_subscriptions(subscriptions)
            user_storage.clear_last_message(chat_id)
            logging.warning(f"User {chat_id} data was removed.")
    except Exception as e:
        logging.error(f"Failed to send message to chat_id={chat_id}: {e}")


def ensure_pipe_exists(pipe_name: str) -> None:
    """Ensure the named pipe exists, creating it if necessary."""
    if not os.path.exists(pipe_name):
        os.mkfifo(pipe_name)
        logging.info(f"Created named pipe at {pipe_name}")


async def listen_pipe(bot) -> None:
    """Listen to the named pipe and send messages to Telegram."""
    ensure_pipe_exists(PIPE_NAME)

    while True:
        with open(PIPE_NAME, "r") as pipe:
            for line in pipe:
                line = line.strip()
                if line:
                    try:
                        first_space = line.find(" ")
                        chat_id = int(line[:first_space])
                        message = line[first_space + 1 :].replace(
                            "\\n", "\n"
                        )  # Convert \\n to actual newlines
                        await send_message(bot, chat_id, message)
                    except Exception as e:
                        logging.error(f"Error while processing message: {e}")


def main() -> None:
    configure_logging()
    load_dotenv()

    token = load_bot_token()
    bot = Bot(token=token)

    logging.info("Starting named pipe listener...")
    asyncio.run(listen_pipe(bot))


if __name__ == "__main__":
    main()
