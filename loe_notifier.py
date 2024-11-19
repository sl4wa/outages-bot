import asyncio
import json
import logging
import re
import sys

from notifier import notifier
from users import user_storage

if __name__ == "__main__":
    logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s - %(levelname)s - %(message)s",
        handlers=[logging.StreamHandler(sys.stdout)]
    )

async def loe_notifier():
    try:
        with open("loe_data.json", "r", encoding="utf-8") as f:
            outages = json.load(f)

        subscriptions = user_storage.load_all_subscriptions()

        for chat_id, subscription in subscriptions.items():
            street_id = subscription.get("street_id")
            street_name = subscription.get("street_name")
            building = subscription.get("building")
            start_date = subscription.get("start_date")
            end_date = subscription.get("end_date")
            comment = subscription.get("comment")

            relevant_outage = next(
                (
                    o
                    for o in outages
                    if str(o["street"]["id"]) == street_id
                    and re.search(rf'\b{building}\b', o["buildingNames"])
                    and (
                        o["dateEvent"] != start_date
                        or o["datePlanIn"] != end_date
                        or o["koment"] != comment
                    )
                ),
                None,
            )

            if relevant_outage:
                try:
                    await notifier.send_message(chat_id, relevant_outage)
                except Exception as e:
                    logging.error(f"Failed to send message to {chat_id}: {e}")
                    return

                subscription['start_date'] = relevant_outage["dateEvent"]
                subscription['end_date'] = relevant_outage["datePlanIn"]
                subscription['comment'] = relevant_outage["koment"]
                user_storage.save_subscription(chat_id, subscription)
                logging.info(f"Notification sent to {chat_id} - {street_name}, {building}")
            else:
                logging.info(f"No relevant outage found for subscription {chat_id} - {street_name}, {building}")
    except (KeyError, ValueError) as e:
        logging.error(f"Error processing outage data: {e}")
    except FileNotFoundError:
        logging.error("loe_data.json file not found!")

if __name__ == "__main__":
    asyncio.run(loe_notifier())
