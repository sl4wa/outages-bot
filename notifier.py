import asyncio
import logging
import re
import sys
from logging.handlers import TimedRotatingFileHandler

from checker import checker
from notifier import notifier
from users import users

LOG_FILE = "notifier.log"


def configure_logging() -> None:
    file_handler = TimedRotatingFileHandler(
        LOG_FILE,
        when="midnight",
        interval=1,
        backupCount=5,
        encoding="utf-8",
        utc=True,
    )
    file_handler.suffix = "%Y-%m-%d.log"

    logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
        datefmt="%Y-%m-%d %H:%M:%S",
        handlers=[file_handler, logging.StreamHandler(sys.stdout)],
    )

    # Suppress excessive logging from HTTP requests library
    httpx_logger = logging.getLogger("httpx")
    httpx_logger.setLevel(logging.WARNING)

    logging.info("Logging is configured.")


async def loe_notifier():
    outages = checker.get_outages()

    subscribed_users = users.all()

    for chat_id, user in subscribed_users.items():
        street_id = user.get("street_id")
        street_name = user.get("street_name")
        building = user.get("building")
        start_date = user.get("start_date")
        end_date = user.get("end_date")
        comment = user.get("comment")

        relevant_outage = next(
            (
                o
                for o in outages
                if str(o["street"]["id"]) == street_id
                and re.search(rf"\b{building}\b", o["buildingNames"])
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

            user["start_date"] = relevant_outage["dateEvent"]
            user["end_date"] = relevant_outage["datePlanIn"]
            user["comment"] = relevant_outage["koment"]
            users.save(chat_id, user)
            logging.info(
                f"Notification sent to {chat_id} - {street_name}, {building}")
        else:
            logging.info(
                f"No relevant outage found for user {chat_id} - {street_name}, {building}"
            )


if __name__ == "__main__":
    configure_logging()
    asyncio.run(loe_notifier())
