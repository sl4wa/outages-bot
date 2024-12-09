import os
from collections.abc import Iterator
from typing import Optional

from .user import User


class UserStorage:
    """File-based implementation for managing users."""

    def __init__(self):
        self.data_directory = "users/data"
        os.makedirs(self.data_directory, exist_ok=True)

    def _get_file_path(self, chat_id: int) -> str:
        """Get the file path for a specific chat ID."""
        return os.path.join(self.data_directory, f"{chat_id}.txt")

    def get(self, chat_id: int) -> Optional[User]:
        """Load a user from a file."""
        file_path = self._get_file_path(chat_id)
        if not os.path.exists(file_path):
            return None

        data = {}
        with open(file_path, encoding="utf-8") as file:
            for line in file:
                if ": " in line:
                    try:
                        key, value = line.strip().split(": ", 1)
                        data[key] = value
                    except ValueError:
                        continue  # Skip lines that do not match the expected format

        return User.from_dict(data)

    def save(self, chat_id: int, user: User) -> None:
        """Save or update a user."""
        file_path = self._get_file_path(chat_id)
        with open(file_path, "w", encoding="utf-8") as file:
            for key, value in user.to_dict().items():
                file.write(f"{key}: {value}\n")

    def remove(self, chat_id: int) -> None:
        """Remove a user."""
        file_path = self._get_file_path(chat_id)
        if os.path.exists(file_path):
            os.remove(file_path)

    def all(self) -> Iterator[tuple[int, User]]:
        """Load all users as a generator."""
        for filename in os.listdir(self.data_directory):
            if filename.endswith(".txt"):
                chat_id = int(filename.replace(".txt", ""))
                user = self.get(chat_id)
                if user:
                    yield chat_id, user
