import logging

from telegram import Update
from telegram.ext import ContextTypes

from users import users


async def handle_stop(update: Update, context: ContextTypes.DEFAULT_TYPE):
    chat_id = update.effective_chat.id
    subscription = users.get(chat_id)

    if subscription:
        users.remove(chat_id)
        await update.message.reply_text(
            "Ви успішно відписалися від сповіщень про відключення електроенергії."
        )
        logging.info(f"User {chat_id} unsubscribed from notifications.")
    else:
        await update.message.reply_text("Ви не маєте активної підписки.")
