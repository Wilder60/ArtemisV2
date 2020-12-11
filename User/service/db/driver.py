"""
This class serves as the interface for any database driver
"""

from service.model import user
from service.model import user_info

class Base():
    def __init__(self, config) -> None:
        pass

    def insert_user(self, user) -> None:
        raise NotImplementedError()

    def get_user(self, user_email) -> user.Login:
        raise NotImplementedError()

    def update_user(self, data) -> None:
        raise NotImplementedError()

    def delete_user(self, data) -> None:
        raise NotImplementedError()

    def get_user_info(self, data) -> user_info.Info:
        raise NotImplementedError()

    def update_user_info(self, data) -> None:
        raise NotImplementedError()