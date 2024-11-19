from .file_users import FileUsers
from .users_interface import UsersInterface

# Global instance of the storage implementation
users: UsersInterface = FileUsers()
