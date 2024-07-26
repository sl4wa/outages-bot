import os
import json
import logging
from telegram import Bot
from dotenv import load_dotenv, find_dotenv

# Setup logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')

# Load environment variables from .env file
dotenv_path = find_dotenv()
if not dotenv_path:
    raise FileNotFoundError("The .env file is missing. Please create a .env file in the project directory with the following content:\n\nTELEGRAM_BOT_TOKEN=your-telegram-bot-token")

load_dotenv(dotenv_path)

# Configuration
TOKEN = os.getenv('TELEGRAM_BOT_TOKEN')
chat_ids_file = 'chat_ids.json'

def load_chat_ids():
    try:
        with open(chat_ids_file, 'r') as file:
            return json.load(file)
    except FileNotFoundError:
        logging.info(f"{chat_ids_file} not found.")
        return []

def list_users():
    chat_ids = load_chat_ids()
    if not chat_ids:
        logging.info("No users are currently subscribed.")
        return

    bot = Bot(token=TOKEN)
    logging.info("Subscribed Users:")
    for chat_id in chat_ids:
        try:
            chat_info = bot.get_chat(chat_id)
            print(f"Chat ID: {chat_info.id}, Username: @{chat_info.username}, First Name: {chat_info.first_name}, Last Name: {chat_info.last_name}")
        except Exception as e:
            logging.error(f"Failed to get info for chat_id {chat_id}: {e}")

if __name__ == '__main__':
    list_users()
