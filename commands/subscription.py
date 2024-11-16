import json
import re

from telegram import Update
from telegram.ext import ContextTypes

from users import user_storage


async def show_subscription(update: Update, context: ContextTypes.DEFAULT_TYPE):
    chat_id = update.effective_chat.id
    subscriptions = user_storage.load_subscriptions()
    last_message = user_storage.load_last_message(chat_id).replace("\\n", "\n")
    last_message = re.sub(r'<.*?>', '', last_message)

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
