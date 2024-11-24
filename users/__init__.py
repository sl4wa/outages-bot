from .user import User
from .users_storage import UsersStorage

__all__ = ["User", "UsersStorage"]

users = UsersStorage()
