import json
import logging

from telegram import Update
from telegram.ext import ContextTypes

from users import user_storage


async def handle_stop(update: Update, context: ContextTypes.DEFAULT_TYPE):
    chat_id = update.effective_chat.id
    subscriptions = user_storage.load_subscriptions()

    if chat_id in subscriptions:
        del subscriptions[chat_id]
        user_storage.save_subscriptions(subscriptions)
        await update.message.reply_text(
            "Ви успішно відписалися від сповіщень про відключення електроенергії."
        )
        logging.info(f"User {chat_id} unsubscribed from notifications.")
    else:
        await update.message.reply_text("Ви не маєте активної підписки.")
