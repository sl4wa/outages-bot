import logging
import sys
from logging.handlers import WatchedFileHandler

from telegram.ext import (
    ApplicationBuilder,
    CommandHandler,
    ConversationHandler,
    MessageHandler,
    filters,
)

from commands.start import building_selection, start, street_selection
from commands.stop import handle_stop
from commands.subscription import show_subscription
from env import load_bot_token

# Constants
LOG_FILE = "bot.log"

# Conversation states
STREET, BUILDING = range(2)


def configure_logging() -> None:
    """Configure logging for the application."""
    logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
        datefmt="%Y-%m-%d %H:%M:%S",
        handlers=[WatchedFileHandler(
            LOG_FILE), logging.StreamHandler(sys.stdout)],
    )

    httpx_logger = logging.getLogger("httpx")
    httpx_logger.setLevel(logging.WARNING)

    logging.info("Logging is configured.")

def setup_bot(token: str):
    """Initialize and set up the Telegram bot with handlers."""
    application = ApplicationBuilder().token(token).build()

    # Conversation handler for /start command
    start_conv_handler = ConversationHandler(
        entry_points=[CommandHandler("start", start)],
        states={
            STREET: [MessageHandler(filters.TEXT & ~filters.COMMAND, street_selection)],
            BUILDING: [
                MessageHandler(filters.TEXT & ~filters.COMMAND,
                               building_selection)
            ],
        },
        fallbacks=[],
        allow_reentry=True,
    )

    # Register handlers
    application.add_handler(start_conv_handler)
    application.add_handler(CommandHandler("subscription", show_subscription))
    application.add_handler(CommandHandler("stop", handle_stop))

    logging.info("Bot handlers are set up.")
    return application


def main() -> None:
    configure_logging()

    application = setup_bot(load_bot_token())

    logging.info("Bot setup completed. Starting polling...")

    # Run the bot until the user presses Ctrl-C
    application.run_polling()
    logging.info("Bot polling has ended.")


if __name__ == "__main__":
    main()
