import unittest
from users.user import User
from outages.outage import Outage

class TestUserOutage(unittest.TestCase):
    def setUp(self):
        # Create users
        self.user1 = User(street_id=12783, street_name="Шевченка Т.", building="271")
        self.user2 = User(street_id=12783, street_name="Шевченка Т.", building="279")
        self.user3 = User(street_id=6458, street_name="Хмельницького Б.", building="294")

        # Create outages
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

        self.outages = [self.outage1, self.outage2]

    def test_get_first_outage(self):
        result1 = self.user1.get_first_outage(self.outages)
        result2 = self.user2.get_first_outage(self.outages)
        result3 = self.user3.get_first_outage(self.outages)

        self.assertEqual(result1, self.outage1)
        self.assertEqual(result2, self.outage1)
        self.assertEqual(result3, self.outage2)

    def test_no_matching_outage(self):
        user = User(street_id=13961, street_name="Залізнична", building="16")
        result = user.get_first_outage(self.outages)
        self.assertIsNone(result)

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

        outages_with_duplicate = [outage_duplicate] + self.outages
        result = self.user1.get_first_outage(outages_with_duplicate)
        self.assertEqual(result, outage_duplicate)

if __name__ == "__main__":
    unittest.main()
