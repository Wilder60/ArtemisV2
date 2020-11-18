from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy import Column, Integer, String

class UserInfo(declarative_base()):
    __tablename__ = "users_info"
    id = Column(Integer, primary_key=True)
    email = Column(String(), nullable=False)
    firstname = Column(String(), nullable=True)
    lastname = Column(String(), nullable=True)
    birthday = Column(String(), nullable=True)
    gender = Column(String(), nullable=True)
    location = Column(String(), nullable=True)
    language = Column(String(), nullable=True)


