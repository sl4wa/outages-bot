import json
import logging
import os
import re
from datetime import datetime

from users import user_storage

PIPE_NAME = "telegram.pipe"

if __name__ == "__main__":
    logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s - %(levelname)s - %(message)s",
        handlers=[logging.StreamHandler()],
    )


def ensure_pipe_exists(pipe_name: str) -> None:
    """Ensure the named pipe exists, creating it if necessary."""
    if not os.path.exists(pipe_name):
        os.mkfifo(pipe_name)
        logging.info(f"Created named pipe at {pipe_name}")


def loe_notifier():
    try:
        with open("loe_data.json", "r", encoding="utf-8") as f:
            outages = json.load(f)
        subscriptions = user_storage.load_subscriptions()

        ensure_pipe_exists(PIPE_NAME)
        with open(PIPE_NAME, "w") as pipe:
            for chat_id, subscription in subscriptions.items():
                relevant_outage = next(
                    (
                        o
                        for o in outages
                        if o["street"]["id"] == subscription["street_id"]
                        and re.search(
                            rf'\b{subscription["building"]}\b', o["buildingNames"]
                        )
                    ),
                    None,
                )
                if relevant_outage:
                    start_time = format_datetime(relevant_outage["dateEvent"])
                    end_time = format_datetime(relevant_outage["datePlanIn"])

                    # Formulate message with actual newlines, then replace with \n for the pipe
                    message = (
                        f"Поточні відключення:\n"
                        f"Місто: {relevant_outage['city']['name']}\n"
                        f"Вулиця: {relevant_outage['street']['name']}\n"
                        f"<b>{start_time} - {end_time}</b>\n"
                        f"Коментар: {relevant_outage['koment']}\n"
                        f"Будинки: {relevant_outage['buildingNames']}"
                    ).replace(
                        "\n", "\\n"
                    )  # Replace actual newlines with \n for transmission

                    last_message = user_storage.load_last_message(chat_id)
                    if message != last_message:
                        # Write message to named pipe
                        pipe.write(f"{chat_id} {message}\n")
                        pipe.flush()  # Ensure the message is written immediately
                        user_storage.save_last_message(chat_id, message)
                        logging.info(
                            f"Notification sent to {chat_id} for subscription: {subscription}"
                        )
                    else:
                        logging.info(
                            f"Outage already notified to {chat_id} for subscription: {subscription}"
                        )
                else:
                    logging.info(
                        f"No relevant outage found for subscription {subscription} of {chat_id}."
                    )
    except (KeyError, ValueError) as e:
        logging.error(f"Error processing outage data: {e}")
    except FileNotFoundError:
        logging.error("loe_data.json file not found!")

def format_datetime(iso_string):
    """Formats the ISO 8601 date string into a readable format."""
    try:
        dt = datetime.fromisoformat(iso_string)
        return dt.strftime("%Y-%m-%d %H:%M")
    except ValueError:
        return iso_string


if __name__ == "__main__":
    loe_notifier()
