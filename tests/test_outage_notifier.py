import asyncio
import logging
import unittest
from unittest.mock import AsyncMock, MagicMock

from telegram.error import Forbidden

from outages import Outage, OutageNotifier, OutageProcessor
from users import User


class TestOutageNotifier(unittest.TestCase):
    def setUp(self):
        self.logger = MagicMock(spec=logging.Logger)
        self.bot = AsyncMock()
        self.user_storage = MagicMock()
        self.outage_processor = OutageProcessor()

        self.outage1 = Outage(
            start_date="2024-11-28T06:47:00+00:00",
            end_date="2024-11-28T10:00:00+00:00",
            city="Львів",
            street_id=12783,
            street="Шевченка Т.",
            building=("271, 273, 273-А, 275, 277, 279, 281, 281-А, 282, 283, 283-А, "
                     "284, 284-А, 285, 285-А, 287, 289, 289-А, 290-А, 291, 291(0083), "
                     "293, 295, 297, 297-А, 297-Б, 308, 313, 316, 316-А, 318, 318-А, "
                     "320, 322, 324, 326, 328, 328-А, 330, 332, 334, 336, 338, 340-А, "
                     "342, 346, 348-А, 350, 350,А, 350-В, 358, 358-А, 360-В"),
            comment="Застосування ГПВ",
        )

        self.outage2 = Outage(
            start_date="2024-11-28T06:47:00+00:00",
            end_date="2024-11-28T10:00:00+00:00",
            city="Львів",
            street_id=6458,
            street="Хмельницького Б.",
            building="294",
            comment="Застосування ГПВ",
        )

        self.outage_processor._outages = [self.outage1, self.outage2]

        self.notifier = OutageNotifier(
            logger=self.logger,
            bot=self.bot,
            user_storage=self.user_storage,
            outage_processor=self.outage_processor,
        )

    def test_outage_already_notified(self):
        user = User(street_id=12783, street_name="Шевченка Т.", building="271")
        user.set_outage(self.outage1)
        self.user_storage.all.return_value = [("chat1", user)]

        asyncio.run(self.notifier.notify())

        self.logger.info.assert_any_call("Outage already notified for user chat1 - Шевченка Т., 271")

    def test_notification_sent(self):
        user = User(street_id=12783, street_name="Шевченка Т.", building="271")
        self.user_storage.all.return_value = [("chat1", user)]

        asyncio.run(self.notifier.notify())

        self.logger.info.assert_any_call("Notification sent to chat1 - Шевченка Т., 271")
        self.bot.send_message.assert_called_once_with(chat_id="chat1", text=self.outage1.format_message(), parse_mode="HTML")

    def test_subscription_removed_for_blocked_user(self):
        user = User(street_id=12783, street_name="Шевченка Т.", building="271")
        self.user_storage.all.return_value = [("chat1", user)]
        self.bot.send_message.side_effect = Forbidden("User blocked the bot")

        asyncio.run(self.notifier.notify())

        self.logger.info.assert_any_call("Subscription removed for blocked user chat1.")
        self.user_storage.remove.assert_called_once_with("chat1")

    def test_no_relevant_outage(self):
        user = User(street_id=99999, street_name="Nonexistent Street", building="1")
        self.user_storage.all.return_value = [("chat1", user)]

        asyncio.run(self.notifier.notify())

        self.logger.info.assert_any_call("No relevant outage found for user chat1 - Nonexistent Street, 1")

if __name__ == "__main__":
    unittest.main()
