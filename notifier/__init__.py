from .notifier_interface import NotifierInterface
from .telegram_notifier import TelegramNotifier

# Global instance of the notifier implementation
notifier: NotifierInterface = TelegramNotifier()
