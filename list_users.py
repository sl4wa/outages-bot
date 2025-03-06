import asyncio
import logging

from telegram import Bot
from telegram.error import TelegramError
from env import load_bot_token

from users import UserStorage

async def list_users():
    bot = Bot(token=load_bot_token())
    print("Subscribed Users:")
    user_storage = UserStorage()
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
