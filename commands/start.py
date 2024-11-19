import json
import logging

from telegram import ReplyKeyboardMarkup, ReplyKeyboardRemove, Update
from telegram.ext import ContextTypes, ConversationHandler

from users import users


def load_streets():
    with open("data/streets.json", "r", encoding="utf-8") as file:
        return json.load(file)


streets = load_streets()
STREET, BUILDING = range(2)


def normalize(text):
    return text.strip().lower()


async def start(update: Update, context: ContextTypes.DEFAULT_TYPE) -> int:
    chat_id = update.effective_chat.id
    subscription = users.get(chat_id)

    if subscription:
        await update.message.reply_text(
            f"Ваша поточна підписка:\nВулиця: {subscription['street_name']}\nБудинок: {subscription['building']}\n\n"
            "Будь ласка, оберіть нову вулицю для оновлення підписки або введіть назву вулиці:"
        )
    else:
        await update.message.reply_text("Будь ласка, введіть назву вулиці:")

    return STREET


async def street_selection(update: Update, context: ContextTypes.DEFAULT_TYPE) -> int:
    query = normalize(update.message.text)
    filtered_streets = [
        street for street in streets if query in normalize(street["name"])
    ]

    if not filtered_streets:
        await update.message.reply_text("Вулицю не знайдено. Спробуйте ще раз.")
        return STREET

    exact_match = next(
        (street for street in filtered_streets if normalize(street["name"]) == query),
        None,
    )
    if exact_match:
        context.user_data["street_name"] = exact_match["name"]
        context.user_data["street_id"] = exact_match["id"]
        await update.message.reply_text(
            text=f"Ви обрали вулицю: {exact_match['name']}\nБудь ласка, введіть номер будинку:",
            reply_markup=ReplyKeyboardRemove(),
        )
        return BUILDING

    keyboard = [[street["name"]] for street in filtered_streets]
    reply_markup = ReplyKeyboardMarkup(
        keyboard, one_time_keyboard=True, resize_keyboard=True
    )
    await update.message.reply_text(
        "Будь ласка, оберіть вулицю:", reply_markup=reply_markup
    )
    return STREET


async def building_selection(update: Update, context: ContextTypes.DEFAULT_TYPE) -> int:
    building = update.message.text

    street_name = context.user_data.get("street_name")
    street_id = context.user_data.get("street_id")
    if not street_name or not street_id:
        await update.message.reply_text(
            "Підписка не завершена. Будь ласка, почніть знову."
        )
        return ConversationHandler.END

    chat_id = update.effective_chat.id
    subscription = {
        "street_id": street_id,
        "street_name": street_name,
        "building": building,
        "start_date": "",
        "end_date": "",
        "comment": ""
    }
    users.save(chat_id, subscription)

    await update.message.reply_text(
        f"Ви підписалися на сповіщення про відключення електроенергії для вулиці {street_name}, будинок {building}.",
        reply_markup=ReplyKeyboardRemove(),
    )

    logging.info(
        f"User {chat_id} subscribed to {street_name} street, building {building}."
    )
    return ConversationHandler.END
