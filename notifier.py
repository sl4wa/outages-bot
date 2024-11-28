import asyncio
import logging
import re
import sys
from logging.handlers import TimedRotatingFileHandler

from telegram.error import Forbidden

from outages import Outage, outages_notifier, outages_reader
from users import User, users

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

    httpx_logger = logging.getLogger("httpx")
    httpx_logger.setLevel(logging.WARNING)

    logging.info("Starting notification script.")


async def notify_user(chat_id: int, user: User, outage: Outage) -> None:
    try:
        await outages_notifier.send_message(chat_id, outage)
        user.start_date = outage.start_date
        user.end_date = outage.end_date
        user.comment = outage.comment
        users.save(chat_id, user)
        logging.info(f"Notification sent to {chat_id} - {user.street_name}, {user.building}")
    except Forbidden:
        users.remove(chat_id)
        logging.info(f"Subscription removed for blocked user {chat_id}.")
    except Exception as e:
        logging.error(f"Failed to send message to {chat_id}: {e}")


async def notifier(outages: list[Outage]) -> None:
    for chat_id, user in users.all():
        # Find the first relevant outage
        outage = next(
            (
                o
                for o in outages
                if o.street_id == user.street_id
                and re.search(rf"\b{re.escape(user.building)}\b", o.building)
            ),
            None,
        )

        if outage:
            # If stored outage matches the first relevant outage, do nothing
            if (
                outage.start_date == user.start_date
                and outage.end_date == user.end_date
                and outage.comment == user.comment
            ):
                logging.info(f"Outage already notified for user {chat_id} - {user.street_name}, {user.building}")
                continue

            # Otherwise, notify about the outage
            await notify_user(chat_id, user, outage)
        else:
            logging.info(f"No relevant outage found for user {chat_id} - {user.street_name}, {user.building}")


def main() -> None:
    configure_logging()
    asyncio.run(notifier(outages_reader.get_outages()))


if __name__ == "__main__":
    main()
