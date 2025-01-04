import logging

from telegram import Update
from telegram.ext import ContextTypes

from users import UserStorage


async def handle_stop(update: Update, context: ContextTypes.DEFAULT_TYPE):
    user_storage = UserStorage()
    chat_id = update.effective_chat.id
    subscription = user_storage.get(chat_id)

    if subscription:
        user_storage.remove(chat_id)
        await update.message.reply_text(
            "Ви успішно відписалися від сповіщень про відключення електроенергії."
        )
        logging.info(f"User {chat_id} unsubscribed from notifications.")
    else:
        await update.message.reply_text("Ви не маєте активної підписки.")
