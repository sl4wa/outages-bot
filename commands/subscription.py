import json

from telegram import Update
from telegram.ext import ContextTypes


async def show_subscription(update: Update, context: ContextTypes.DEFAULT_TYPE):
    chat_id = str(update.effective_chat.id)
    subscriptions = load_subscriptions()
    last_message = load_last_message(chat_id)

    if chat_id in subscriptions:
        current_sub = subscriptions[chat_id]
        message = f"Ваша поточна підписка:\nВулиця: {current_sub['street_name']}\nБудинок: {current_sub['building']}"
        if last_message:
            message += f"\n\nОстаннє відправлене повідомлення:\n{last_message}"
    else:
        message = "Ви не маєте активної підписки."

    if update.message:
        await update.message.reply_text(message)
    elif update.callback_query:
        await update.callback_query.answer()
        await update.callback_query.message.reply_text(message)
    else:
        await context.bot.send_message(chat_id=chat_id, text=message)


def load_subscriptions():
    try:
        with open("subscriptions.json", "r", encoding="utf-8") as file:
            return json.load(file)
    except FileNotFoundError:
        return {}


def load_last_message(chat_id):
    try:
        with open("last_messages.json", "r", encoding="utf-8") as file:
            data = json.load(file)
            return data.get(chat_id, "")
    except FileNotFoundError:
        return ""
