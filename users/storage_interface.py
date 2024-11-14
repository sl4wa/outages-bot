# users/storage_interface.py

from abc import ABC, abstractmethod
from typing import Dict, Any

class UserStorage(ABC):
    """Interface for user data storage."""

    @abstractmethod
    def load_subscriptions(self) -> Dict[int, Any]:
        pass

    @abstractmethod
    def save_subscriptions(self, subscriptions: Dict[int, Any]) -> None:
        pass

    @abstractmethod
    def load_last_message(self, chat_id: int) -> str:
        pass

    @abstractmethod
    def save_last_message(self, chat_id: int, message: str) -> None:
        pass

    @abstractmethod
    def clear_last_message(self, chat_id: int) -> None:
        pass
