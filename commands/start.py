from telegram import ReplyKeyboardMarkup, ReplyKeyboardRemove, Update
from telegram.ext import CallbackContext, ConversationHandler
import json
import logging

def load_streets():
    with open('streets.json', 'r', encoding='utf-8') as file:
        data = json.load(file)
        return data['hydra:member']

streets = load_streets()
STREET, BUILDING = range(2)

def normalize(text):
    return text.strip().lower()

def start(update: Update, context: CallbackContext) -> int:
    chat_id = update.effective_chat.id
    subscriptions = load_subscriptions()

    if str(chat_id) in subscriptions:
        current_sub = subscriptions[str(chat_id)]
        update.message.reply_text(
            f"Ваша поточна підписка:\nВулиця: {current_sub['street_name']}\nБудинок: {current_sub['building']}\n\n"
            "Будь ласка, оберіть нову вулицю для оновлення підписки або введіть назву вулиці:"
        )
    else:
        update.message.reply_text("Будь ласка, введіть назву вулиці:")

    return STREET

def street_selection(update: Update, context: CallbackContext) -> int:
    query = normalize(update.message.text)
    filtered_streets = [street for street in streets if query in normalize(street['name'])]

    if not filtered_streets:
        update.message.reply_text('Вулицю не знайдено. Спробуйте ще раз.')
        return STREET

    exact_match = next((street for street in filtered_streets if normalize(street['name']) == query), None)
    if exact_match:
        context.user_data['street_name'] = exact_match['name']
        context.user_data['street_id'] = exact_match['id']
        update.message.reply_text(
            text=f"Ви обрали вулицю: {exact_match['name']}\nБудь ласка, введіть номер будинку:",
            reply_markup=ReplyKeyboardRemove()
        )
        return BUILDING

    keyboard = [[street['name']] for street in filtered_streets]
    reply_markup = ReplyKeyboardMarkup(keyboard, one_time_keyboard=True, resize_keyboard=True)
    update.message.reply_text('Будь ласка, оберіть вулицю:', reply_markup=reply_markup)
    return STREET

def building_selection(update: Update, context: CallbackContext) -> int:
    building = update.message.text

    street_name = context.user_data.get('street_name')
    street_id = context.user_data.get('street_id')
    if not street_name or not street_id:
        update.message.reply_text('Підписка не завершена. Будь ласка, почніть знову.')
        return ConversationHandler.END

    chat_id = str(update.effective_chat.id)
    subscriptions = load_subscriptions()
    subscriptions[chat_id] = {"street_id": street_id, "street_name": street_name, "building": building}
    save_subscriptions(subscriptions)
    clear_last_message(chat_id)

    update.message.reply_text(
        f"Ви підписалися на сповіщення про відключення електроенергії для вулиці {street_name}, будинок {building}.",
        reply_markup=ReplyKeyboardRemove()
    )

    logging.info(f"User {chat_id} subscribed to {street_name} street, building {building}.")
    return ConversationHandler.END

def load_subscriptions():
    try:
        with open('subscriptions.json', 'r', encoding='utf-8') as file:
            return json.load(file)
    except FileNotFoundError:
        return {}

def save_subscriptions(subscriptions):
    with open('subscriptions.json', 'w', encoding='utf-8') as file:
        json.dump({str(k): v for k, v in subscriptions.items()}, file, ensure_ascii=False, indent=4)

def clear_last_message(chat_id):
    try:
        with open('last_messages.json', 'r', encoding='utf-8') as file:
            data = json.load(file)
    except FileNotFoundError:
        data = {}
    data.pop(str(chat_id), None)
    with open('last_messages.json', 'w', encoding='utf-8') as file:
        json.dump(data, file, ensure_ascii=False, indent=4)
