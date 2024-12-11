import asyncio
import logging
import os

from dotenv import find_dotenv, load_dotenv
from telegram import Bot
from telegram.error import TelegramError

from users import user_storage

# Load environment variables from .env file
dotenv_path = find_dotenv()
if not dotenv_path:
    raise FileNotFoundError(
        "The .env file is missing. Please create a .env file in the project directory with the following content:\n\nTELEGRAM_BOT_TOKEN=your-telegram-bot-token"
    )

load_dotenv(dotenv_path)

# Configuration
TOKEN = os.getenv("TELEGRAM_BOT_TOKEN")

async def list_users():
    bot = Bot(token=TOKEN)
    print("Subscribed Users:")
    users = user_storage.all()
    user_count = 0
    for chat_id, user in users:
        try:
            chat_info = await bot.get_chat(chat_id)
            print(
                f"Chat ID: {chat_info.id}, Username: @{chat_info.username}, First Name: {chat_info.first_name}, "
                f"Last Name: {chat_info.last_name}, "
                f"Street Name: {user.street_name}, Building: {user.building}"
            )
            user_count += 1
        except TelegramError as e:
            logging.error(f"Failed to get info for chat_id {chat_id}: {e}")

    print(f"Total Users: {user_count}")


if __name__ == "__main__":
    asyncio.run(list_users())
