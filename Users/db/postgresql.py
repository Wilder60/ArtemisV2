from sqlalchemy import create_engine
from typing import Dict, Any

class postgresql():
    
    def __init__(self, config: Dict[Any]):
        engine = create_engine("sqllite:///:memory:", echo=True)
        pass

    def get_user(self, user_id):
        pass

    def get_user_info(self, user_id):
        pass

    def insert_user(self, data):
        pass

    def update_user(self, data):
        pass

    def delete_user(self, id):
        pass