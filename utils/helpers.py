import json
import datetime

last_messages_file = 'last_messages.json'

def load_last_message(chat_id):
    try:
        with open(last_messages_file, 'r', encoding='utf-8') as file:
            data = json.load(file)
            return data.get(str(chat_id), '')
    except FileNotFoundError:
        return ''

def save_last_message(chat_id, message):
    try:
        with open(last_messages_file, 'r', encoding='utf-8') as file:
            data = json.load(file)
    except FileNotFoundError:
        data = {}
    data[str(chat_id)] = message
    with open(last_messages_file, 'w', encoding='utf-8') as file:
        json.dump(data, file, ensure_ascii=False, indent=4)

def format_datetime(date_string):
    dt = datetime.datetime.fromisoformat(date_string)
    return dt.strftime('%Y-%m-%d %H:%M')
