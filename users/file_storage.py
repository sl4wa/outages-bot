import os
from typing import Dict, Optional

class FileUserStorage:
    """File-based implementation of the UserStorage interface using key-value storage."""

    def __init__(self):
        self.data_directory = "users/data"
        os.makedirs(self.data_directory, exist_ok=True)

    def _get_file_path(self, chat_id: int) -> str:
        """Get the file path for a specific chat ID."""
        return os.path.join(self.data_directory, f"{chat_id}.txt")

    def load_subscription(self, chat_id: int) -> Optional[Dict[str, str]]:
        """Load a subscription from a file."""
        file_path = self._get_file_path(chat_id)
        if not os.path.exists(file_path):
            return None

        subscription = {}
        with open(file_path, "r", encoding="utf-8") as file:
            for line in file:
                if ": " in line:
                    try:
                        key, value = line.strip().split(": ", 1)
                        subscription[key] = value
                    except ValueError:
                        continue  # Handle lines that do not match the expected format
        return subscription

    def save_subscription(self, chat_id: int, subscription: Dict[str, str]) -> None:
        """Save or update a subscription."""
        file_path = self._get_file_path(chat_id)
        with open(file_path, "w", encoding="utf-8") as file:
            for key, value in subscription.items():
                file.write(f"{key}: {value}\n")

    def remove_subscription(self, chat_id: int) -> None:
        """Remove a subscription."""
        file_path = self._get_file_path(chat_id)
        if os.path.exists(file_path):
            os.remove(file_path)

    def load_all_subscriptions(self) -> Dict[int, Dict[str, str]]:
        """Load all subscriptions."""
        subscriptions = {}
        for filename in os.listdir(self.data_directory):
            if filename.endswith(".txt"):
                chat_id = int(filename.replace(".txt", ""))
                subscription = self.load_subscription(chat_id)
                if subscription:
                    subscriptions[chat_id] = subscription
        return subscriptions
