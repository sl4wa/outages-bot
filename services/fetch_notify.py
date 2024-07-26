import requests
import logging
from telegram import Bot, ParseMode
from utils.helpers import save_last_message, format_datetime, load_last_message
import os
import re
from dotenv import load_dotenv
import json

load_dotenv()

api_url = 'https://power-api.loe.lviv.ua/api/pw_accidents?pagination=false&otg.id=28&city.id=693'
TOKEN = os.getenv('TELEGRAM_BOT_TOKEN')

def fetch_and_notify():
    logging.info("Attempting to fetch outage data from API...")
    response = requests.get(api_url)
    
    if response.status_code == 200:
        try:
            data = response.json()
            outages = data.get('hydra:member', [])
            subscriptions = load_subscriptions()

            bot = Bot(token=TOKEN)
            for chat_id, subscription in subscriptions.items():
                relevant_outage = next((o for o in outages if o['street']['id'] == subscription['street_id'] and re.search(rf'\b{subscription["building"]}\b', o['buildingNames'])), None)
                if relevant_outage:
                    start_time = format_datetime(relevant_outage['dateEvent'])
                    end_time = format_datetime(relevant_outage['datePlanIn'])

                    message = (
                        f"Поточні відключення:\n"
                        f"Місто: {relevant_outage['city']['name']}\n"
                        f"Вулиця: {relevant_outage['street']['name']}\n"
                        f"*{start_time} - {end_time}*\n"
                        f"Коментар: {relevant_outage['koment']}\n"
                        f"Будинки: {relevant_outage['buildingNames']}"
                    )

                    last_message = load_last_message(chat_id)
                    if message != last_message:
                        try:
                            bot.send_message(chat_id=chat_id, text=message, parse_mode=ParseMode.MARKDOWN)
                            save_last_message(chat_id, message)
                            logging.info(f"Notification sent to {chat_id} for subscription: {subscription}")
                        except Exception as e:
                            logging.error(f"Failed to send message to chat_id {chat_id}: {e}")
                    else:
                        logging.info(f"Outage already notified to {chat_id} for subscription: {subscription}")
                else:
                    logging.info(f"No relevant outage found for subscription {subscription} of {chat_id}.")
        except (KeyError, ValueError) as e:
            logging.error(f"Error processing API response: {e}")
    else:
        logging.error(f"Failed to fetch data: HTTP {response.status_code}")

def load_subscriptions():
    try:
        with open('subscriptions.json', 'r', encoding='utf-8') as file:
            return json.load(file)
    except FileNotFoundError:
        return {}
