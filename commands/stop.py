import logging
from telegram import Update
from telegram.ext import CallbackContext
import json

def load_subscriptions():
    try:
        with open('subscriptions.json', 'r', encoding='utf-8') as file:
            return json.load(file)
    except FileNotFoundError:
        return {}

def save_subscriptions(subscriptions):
    with open('subscriptions.json', 'w', encoding='utf-8') as file:
        json.dump(subscriptions, file, ensure_ascii=False, indent=4)

def handle_stop(update: Update, context: CallbackContext):
    chat_id = str(update.effective_chat.id)
    subscriptions = load_subscriptions()

    if chat_id in subscriptions:
        del subscriptions[chat_id]
        save_subscriptions(subscriptions)
        update.message.reply_text("Ви успішно відписалися від сповіщень про відключення електроенергії.")
        logging.info(f"User {chat_id} unsubscribed from notifications.")
    else:
        update.message.reply_text("Ви не маєте активної підписки.")
