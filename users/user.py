from dataclasses import asdict, dataclass
from typing import Optional, Self


@dataclass
class User:
    """Data class representing a user."""
    street_id: int
    street_name: str
    building: str
    start_date: Optional[str] = None
    end_date: Optional[str] = None
    comment: Optional[str] = None

    @staticmethod
    def from_dict(data: dict[str, str]) -> "User":
        """Create a User from a dictionary."""
        return User(
            street_id=int(data.get("street_id", 0)),
            street_name=data.get("street_name", ""),
            building=data.get("building", ""),
            start_date=data.get("start_date"),
            end_date=data.get("end_date"),
            comment=data.get("comment")
        )

    def to_dict(self) -> dict[str, str]:
        """Convert the User to a dictionary."""
        return {key: str(value) for key, value in asdict(self).items() if value is not None}

    def is_notified(self, outage: "Outage") -> bool:
        """Check if the outage is already notified for the user."""
        return (
            outage.start_date == self.start_date
            and outage.end_date == self.end_date
            and outage.comment == self.comment
        )

    def set_outage(self, outage: "Outage") -> Self:
        """Update the user's outage information."""
        self.start_date = outage.start_date
        self.end_date = outage.end_date
        self.comment = outage.comment

        return self