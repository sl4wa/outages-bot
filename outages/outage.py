from dataclasses import dataclass


@dataclass
class Outage:
    start_date: str
    end_date: str
    city: str
    street_id: int
    street: str
    building: str
    comment: str