import unittest

from outages import Outage, OutageProcessor
from users import User


class TestOutageProcessor(unittest.TestCase):
    def setUp(self):
        self.outage_processor = OutageProcessor()

        outage1 = Outage(
            start_date="2024-11-28T06:47:00+00:00",
            end_date="2024-11-28T10:00:00+00:00",
            city="Львів",
            street_id=12783,
            street="Шевченка Т.",
            building=(
                "271, 273, 273-А, 275, 277, 279, 281, 281-А, 282, 283, 283-А, "
                "284, 284-А, 285, 285-А, 287, 289, 289-А, 290-А, 291, 291(0083), "
                "293, 295, 297, 297-А, 297-Б, 308, 313, 316, 316-А, 318, 318-А, "
                "320, 322, 324, 326, 328, 328-А, 330, 332, 334, 336, 338, 340-А, "
                "342, 346, 348-А, 350, 350,А, 350-В, 358, 358-А, 360-В"
            ),
            comment="Застосування ГПВ",
        )

        outage2 = Outage(
            start_date="2024-11-28T06:47:00+00:00",
            end_date="2024-11-28T10:00:00+00:00",
            city="Львів",
            street_id=6458,
            street="Хмельницького Б.",
            building="294",
            comment="Застосування ГПВ",
        )

        self.outage_processor._outages = [outage1, outage2]

    def test_get_user_outage(self):
        user1 = User(street_id=12783, street_name="Шевченка Т.", building="271")
        user2 = User(street_id=12783, street_name="Шевченка Т.", building="279")
        user3 = User(street_id=6458, street_name="Хмельницького Б.", building="294")

        self.assertEqual(self.outage_processor.get_user_outage(user1).street_id, 12783)
        self.assertEqual(self.outage_processor.get_user_outage(user2).street_id, 12783)
        self.assertEqual(self.outage_processor.get_user_outage(user3).street_id, 6458)

    def test_no_matching_outage(self):
        user = User(street_id=13961, street_name="Залізнична", building="16")
        self.assertIsNone(self.outage_processor.get_user_outage(user))

    def test_multiple_outages_same_building(self):
        outage_duplicate = Outage(
            start_date="2024-11-28T06:47:00+00:00",
            end_date="2024-11-28T10:00:00+00:00",
            city="Львів",
            street_id=12783,
            street="Шевченка Т.",
            building="271",
            comment="Застосування ГАВ",
        )

        self.outage_processor._outages.insert(0, outage_duplicate)

        user = User(street_id=12783, street_name="Шевченка Т.", building="271")
        self.assertEqual(self.outage_processor.get_user_outage(user).comment, "Застосування ГАВ")

if __name__ == "__main__":
    unittest.main()
