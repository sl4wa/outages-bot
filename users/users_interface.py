from abc import ABC, abstractmethod
from typing import Dict, Optional

class UsersInterface(ABC):
    """Interface for user data storage."""

    @abstractmethod
    def get(self, chat_id: int) -> Optional[Dict[str, str]]:
        """Load a specific subscription."""
        pass

    @abstractmethod
    def save(self, chat_id: int, subscription: Dict[str, str]) -> None:
        """Save or update a subscription."""
        pass

    @abstractmethod
    def remove(self, chat_id: int) -> None:
        """Remove a subscription."""
        pass

    @abstractmethod
    def all(self) -> Dict[int, Dict[str, str]]:
        """Load all subscriptions."""
        pass
