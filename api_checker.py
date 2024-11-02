import os
import logging
import requests
import re
import json
from datetime import datetime

# Constants
PIPE_NAME = 'telegram.pipe'
API_URL = 'https://power-api.loe.lviv.ua/api/pw_accidents?pagination=false&otg.id=28&city.id=693'
LOG_FILE = 'api_checker.log'

# Set up logging to log both to a file and stdout
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler(LOG_FILE),
        logging.StreamHandler()
    ]
)

HEADERS = {
    'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3',
    'Accept': 'application/json, text/plain, */*',
    'Connection': 'keep-alive',
    'Accept-Language': 'en-US,en;q=0.9',
}

def ensure_pipe_exists(pipe_name: str) -> None:
    """Ensure the named pipe exists, creating it if necessary."""
    if not os.path.exists(pipe_name):
        os.mkfifo(pipe_name)
        logging.info(f"Created named pipe at {pipe_name}")

def notify():
    logging.info("Attempting to fetch outage data from API...")
    response = requests.get(API_URL, headers=HEADERS)
    
    if response.status_code == 200:
        try:
            data = response.json()
            outages = data.get('hydra:member', [])
            subscriptions = load_subscriptions()

            ensure_pipe_exists(PIPE_NAME)
            with open(PIPE_NAME, 'w') as pipe:
                for chat_id, subscription in subscriptions.items():
                    relevant_outage = next(
                        (o for o in outages if o['street']['id'] == subscription['street_id'] and re.search(rf'\b{subscription["building"]}\b', o['buildingNames'])),
                        None
                    )
                    if relevant_outage:
                        start_time = format_datetime(relevant_outage['dateEvent'])
                        end_time = format_datetime(relevant_outage['datePlanIn'])

                        # Formulate message with actual newlines, then replace with \n for the pipe
                        message = (
                            f"Поточні відключення:\n"
                            f"Місто: {relevant_outage['city']['name']}\n"
                            f"Вулиця: {relevant_outage['street']['name']}\n"
                            f"<b>{start_time} - {end_time}</b>\n"
                            f"Коментар: {relevant_outage['koment']}\n"
                            f"Будинки: {relevant_outage['buildingNames']}"
                        ).replace('\n', '\\n')  # Replace actual newlines with \n for transmission

                        last_message = load_last_message(chat_id)
                        if message != last_message:
                            # Write message to named pipe
                            pipe.write(f"{chat_id} {message}\n")
                            pipe.flush()  # Ensure the message is written immediately
                            save_last_message(chat_id, message)
                            logging.info(f"Notification sent to {chat_id} for subscription: {subscription}")
                        else:
                            logging.info(f"Outage already notified to {chat_id} for subscription: {subscription}")
                    else:
                        logging.info(f"No relevant outage found for subscription {subscription} of {chat_id}.")
        except (KeyError, ValueError) as e:
            logging.error(f"Error processing API response: {e}")
    else:
        logging.error(f"Failed to fetch data: HTTP {response.status_code}")

def load_subscriptions():
    """Loads the subscription data from a JSON file."""
    try:
        with open('subscriptions.json', 'r', encoding='utf-8') as file:
            return json.load(file)
    except FileNotFoundError:
        logging.error("subscriptions.json file not found!")
        return {}

def save_last_message(chat_id, message):
    """Saves the last message sent to each chat_id."""
    try:
        with open('last_messages.json', 'r+', encoding='utf-8') as file:
            last_messages = json.load(file)
            last_messages[chat_id] = message
            file.seek(0)
            json.dump(last_messages, file, ensure_ascii=False, indent=4)
            file.truncate()
    except FileNotFoundError:
        with open('last_messages.json', 'w', encoding='utf-8') as file:
            json.dump({chat_id: message}, file, ensure_ascii=False, indent=4)

def load_last_message(chat_id):
    """Loads the last message sent to a chat_id."""
    try:
        with open('last_messages.json', 'r', encoding='utf-8') as file:
            last_messages = json.load(file)
            return last_messages.get(chat_id, "")
    except FileNotFoundError:
        return ""

def format_datetime(iso_string):
    """Formats the ISO 8601 date string into a readable format."""
    try:
        dt = datetime.fromisoformat(iso_string)
        return dt.strftime('%Y-%m-%d %H:%M')
    except ValueError:
        return iso_string

if __name__ == "__main__":
    notify()
