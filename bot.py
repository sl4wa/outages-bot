import warnings

# Suppress specific warnings about urllib3 before importing telegram modules
warnings.filterwarnings(
    "ignore",
    category=UserWarning,
    message=(
        "python-telegram-bot is using upstream urllib3. This is allowed "
        "but not supported by python-telegram-bot maintainers."
    )
)

import logging
import os
import signal
import sys

from dotenv import load_dotenv
from logging.handlers import WatchedFileHandler
from telegram.ext import (
    Updater,
    CommandHandler,
    MessageHandler,
    Filters,
    ConversationHandler
)

# Local imports
from commands.start import start, street_selection, building_selection
from commands.subscription import show_subscription
from commands.stop import handle_stop

# Constants
LOG_FILE = 'bot.log'
TELEGRAM_TOKEN_ENV = 'TELEGRAM_BOT_TOKEN'

# Conversation states
STREET, BUILDING = range(2)

def configure_logging() -> None:
    """Configure logging for the application."""
    logging.basicConfig(
        level=logging.INFO,
        format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
        datefmt='%Y-%m-%d %H:%M:%S',
        handlers=[
            WatchedFileHandler(LOG_FILE),
            logging.StreamHandler(sys.stdout)
        ]
    )
    logging.info("Logging is configured.")


def shutdown(signum: int, frame) -> None:
    """Gracefully shutdown the application."""
    logging.info("Received shutdown signal. Shutting down gracefully...")
    sys.exit(0)


def setup_signal_handlers() -> None:
    """Register signal handlers for graceful shutdown."""
    signal.signal(signal.SIGINT, shutdown)
    signal.signal(signal.SIGTERM, shutdown)
    logging.info("Signal handlers registered.")


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


def setup_bot(token: str) -> Updater:
    """Initialize and set up the Telegram bot with handlers."""
    updater = Updater(token, use_context=True)
    dispatcher = updater.dispatcher

    # Conversation handler for /start command
    start_conv_handler = ConversationHandler(
        entry_points=[CommandHandler('start', start)],
        states={
            STREET: [MessageHandler(Filters.text & ~Filters.command, street_selection)],
            BUILDING: [MessageHandler(Filters.text & ~Filters.command, building_selection)]
        },
        fallbacks=[],
        allow_reentry=True
    )

    # Register handlers
    dispatcher.add_handler(start_conv_handler)
    dispatcher.add_handler(CommandHandler('subscription', show_subscription))
    dispatcher.add_handler(CommandHandler('stop', handle_stop))

    logging.info("Bot handlers are set up.")
    return updater


def main() -> None:
    configure_logging()
    load_dotenv()

    token = load_bot_token()
    setup_signal_handlers()

    updater = setup_bot(token)

    logging.info("Bot setup completed. Starting polling...")
    updater.start_polling()
    logging.info("Bot is now polling. Press Ctrl+C to stop.")

    updater.idle()
    logging.info("Bot polling has ended.")


if __name__ == '__main__':
    main()
