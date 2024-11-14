# users/file_storage.py

import json
from typing import Dict, Any
from .storage_interface import UserStorage

class FileUserStorage(UserStorage):
    """File-based implementation of the UserStorage interface."""

    def __init__(self):
        self.subscriptions_file = "subscriptions.json"
        self.last_messages_file = "last_messages.json"

    def load_subscriptions(self) -> Dict[int, Any]:
        """Load subscriptions from a JSON file."""
        try:
            with open(self.subscriptions_file, "r", encoding="utf-8") as file:
                data = json.load(file)
                # Automatically handle chat_id as integers internally
                return {int(chat_id): info for chat_id, info in data.items()}
        except FileNotFoundError:
            return {}

    def save_subscriptions(self, subscriptions: Dict[int, Any]) -> None:
        """Save subscriptions to a JSON file."""
        with open(self.subscriptions_file, "w", encoding="utf-8") as file:
            json.dump(subscriptions, file, ensure_ascii=False, indent=4)

    def load_last_message(self, chat_id: int) -> str:
        """Load the last message for a specific chat ID."""
        try:
            with open(self.last_messages_file, "r", encoding="utf-8") as file:
                data = json.load(file)
                return data.get(str(chat_id), "")
        except FileNotFoundError:
            return ""

    def save_last_message(self, chat_id: int, message: str) -> None:
        """Save the last message for a specific chat ID."""
        try:
            with open(self.last_messages_file, "r+", encoding="utf-8") as file:
                last_messages = json.load(file)
                last_messages[chat_id] = message
                file.seek(0)
                json.dump(last_messages, file, ensure_ascii=False, indent=4)
                file.truncate()
        except FileNotFoundError:
            with open(self.last_messages_file, "w", encoding="utf-8") as file:
                json.dump({chat_id: message}, file, ensure_ascii=False, indent=4)

    def clear_last_message(self, chat_id: int) -> None:
        """Clear the last message for a specific chat ID."""
        try:
            with open(self.last_messages_file, "r", encoding="utf-8") as file:
                data = json.load(file)
        except FileNotFoundError:
            data = {}
        data.pop(chat_id, None)
        with open(self.last_messages_file, "w", encoding="utf-8") as file:
            json.dump(data, file, ensure_ascii=False, indent=4)
