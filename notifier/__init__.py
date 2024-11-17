from .telegram_notifier import TelegramNotifier
from .notifier_interface import NotifierInterface

# Global instance of the notifier implementation
notifier: NotifierInterface = TelegramNotifier()
