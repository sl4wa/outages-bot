from telegram import Update
from telegram.ext import ContextTypes

from users import user_storage


async def show_subscription(update: Update, context: ContextTypes.DEFAULT_TYPE):
    chat_id = update.effective_chat.id
    subscription = user_storage.get(chat_id)

    if subscription:
        message = f"Ваша поточна підписка:\nВулиця: {subscription.street_name}\nБудинок: {subscription.building}"
    else:
        message = "Ви не маєте активної підписки."

    if update.message:
        await update.message.reply_text(message)
    elif update.callback_query:
        await update.callback_query.answer()
        await update.callback_query.message.reply_text(message)
    else:
        await context.bot.send_message(chat_id=chat_id, text=message)
