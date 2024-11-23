from dataclasses import dataclass
from typing import Optional, Dict, List

@dataclass
class Outage:
    start_date: str
    end_date: str
    city: str
    street_id: int
    street: str
    building: str
    comment: str