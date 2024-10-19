import logging
import os
import signal
import sys
import warnings
from logging.handlers import WatchedFileHandler

# Suppress specific warning about urllib3
warnings.filterwarnings("ignore", category=UserWarning, message="python-telegram-bot is using upstream urllib3. This is allowed but not supported by python-telegram-bot maintainers.")

from dotenv import load_dotenv
from bot_setup import setup_bot, get_scheduler

# Configure logging
log_file = 'bot.log'
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    datefmt='%Y-%m-%d %H:%M:%S',
    handlers=[
        WatchedFileHandler(log_file),
        logging.StreamHandler()
    ]
)

def shutdown(signum, frame):
    logging.info("Shutting down gracefully...")
    scheduler = get_scheduler()
    scheduler.shutdown(wait=False)
    sys.exit(0)

if __name__ == '__main__':
    # Load environment variables from .env file
    load_dotenv()

    token = os.getenv('TELEGRAM_BOT_TOKEN')
    if not token:
        raise ValueError("No token provided! Please add the TELEGRAM_BOT_TOKEN to the .env file as follows:\n\nTELEGRAM_BOT_TOKEN=your-telegram-bot-token")

    # Set up signal handling
    signal.signal(signal.SIGINT, shutdown)
    signal.signal(signal.SIGTERM, shutdown)

    # Set up the bot
    setup_bot(token)
