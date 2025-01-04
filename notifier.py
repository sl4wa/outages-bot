import asyncio
import logging
import sys
from logging.handlers import TimedRotatingFileHandler

from outages import OutageNotifier

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

async def main() -> None:
    logger = logging.getLogger("OutageNotifier")
    outage_notifier = OutageNotifier(logger=logger)
    await outage_notifier.notify()

if __name__ == "__main__":
    configure_logging()
    asyncio.run(main())
