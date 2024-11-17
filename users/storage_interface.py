from abc import ABC, abstractmethod
from typing import Dict, Optional

class UserStorage(ABC):
    """Interface for user data storage."""

    @abstractmethod
    def load_subscription(self, chat_id: int) -> Optional[Dict[str, str]]:
        """Load a specific subscription."""
        pass

    @abstractmethod
    def save_subscription(self, chat_id: int, subscription: Dict[str, str]) -> None:
        """Save or update a subscription."""
        pass

    @abstractmethod
    def remove_subscription(self, chat_id: int) -> None:
        """Remove a subscription."""
        pass

    @abstractmethod
    def load_all_subscriptions(self) -> Dict[int, Dict[str, str]]:
        """Load all subscriptions."""
        pass
