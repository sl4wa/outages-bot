# users/__init__.py

from .file_storage import FileUserStorage
from .storage_interface import UserStorage

# Global instance of the storage implementation
user_storage: UserStorage = FileUserStorage()
