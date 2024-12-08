from dataclasses import dataclass
from datetime import datetime


@dataclass
class Outage:
    start_date: str
    end_date: str
    city: str
    street_id: int
    street: str
    building: str
    comment: str

    def format_date(self, iso_string: str) -> str:
        """Formats the ISO 8601 date string into a readable format."""
        try:
            dt = datetime.fromisoformat(iso_string)
            return dt.strftime("%Y-%m-%d %H:%M")
        except ValueError:
            return iso_string

    def format_message(self) -> str:
        """Formats the outage details into a message."""
        start = self.format_date(self.start_date)
        end = self.format_date(self.end_date)
        message = (
            f"Поточні відключення:\n"
            f"Місто: {self.city}\n"
            f"Вулиця: {self.street}\n"
            f"<b>{start} - {end}</b>\n"
            f"Коментар: {self.comment}\n"
            f"Будинки: {self.building}"
        )
        return message
