import asyncio
import json
import logging
import os

from dotenv import find_dotenv, load_dotenv
from telegram import Bot
from telegram.error import TelegramError

# Setup logging
logging.basicConfig(
    level=logging.INFO, format="%(asctime)s - %(name)s - %(levelname)s - %(message)s"
)

# Load environment variables from .env file
dotenv_path = find_dotenv()
if not dotenv_path:
    raise FileNotFoundError(
        "The .env file is missing. Please create a .env file in the project directory with the following content:\n\nTELEGRAM_BOT_TOKEN=your-telegram-bot-token"
    )

load_dotenv(dotenv_path)

# Configuration
TOKEN = os.getenv("TELEGRAM_BOT_TOKEN")
subscriptions_file = "subscriptions.json"


def load_chat_data():
    try:
        with open(subscriptions_file, "r") as file:
            return json.load(file)
    except FileNotFoundError:
        logging.info(f"{subscriptions_file} not found.")
        return {}


async def list_users():
    chat_data = load_chat_data()
    if not chat_data:
        logging.info("No users are currently subscribed.")
        return

    bot = Bot(token=TOKEN)
    logging.info("Subscribed Users:")
    for chat_id, info in chat_data.items():
        try:
            chat_info = await bot.get_chat(chat_id)
            print(
                f"Chat ID: {chat_info.id}, Username: @{chat_info.username}, First Name: {chat_info.first_name}, "
                f"Last Name: {chat_info.last_name}, Street ID: {info['street_id']}, "
                f"Street Name: {info['street_name']}, Building: {info['building']}"
            )
        except TelegramError as e:
            logging.error(f"Failed to get info for chat_id {chat_id}: {e}")


if __name__ == "__main__":
    asyncio.run(list_users())
