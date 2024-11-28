from datetime import datetime

from .outage import Outage


class OutagesFormatter:

    def _format_datetime(self, iso_string: str) -> str:
        """Formats the ISO 8601 date string into a readable format."""
        try:
            dt = datetime.fromisoformat(iso_string)
            return dt.strftime("%Y-%m-%d %H:%M")
        except ValueError:
            return iso_string

    def format_message(self, outage: Outage) -> str:
        """Send a message to the specified Telegram chat ID."""
        start = self._format_datetime(outage.start_date)
        end = self._format_datetime(outage.end_date)
        message = (
            f"Поточні відключення:\n"
            f"Місто: {outage.city}\n"
            f"Вулиця: {outage.street}\n"
            f"<b>{start} - {end}</b>\n"
            f"Коментар: {outage.comment}\n"
            f"Будинки: {outage.building}"
        )

        return message
