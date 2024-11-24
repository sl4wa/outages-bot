from .outage import Outage
from .outages_notifier import OutagesNotifier
from .outages_reader import OutagesReader

__all__ = ["Outage", "OutagesNotifier", "OutagesReader"]

outages_reader = OutagesReader()
outages_notifier = OutagesNotifier()
