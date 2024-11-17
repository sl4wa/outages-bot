from abc import ABC, abstractmethod

class NotifierInterface(ABC):
    """Interface for different types of notifiers."""

    @abstractmethod
    async def send_message(self, chat_id: int, message: str) -> None:
        """Send a message to a given chat ID."""
        pass
