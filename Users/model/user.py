from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy import Column, Integer, String

class User(declarative_base):
    __tablename__ = "user"
    id = Column(Integer, primary_key=True)
    email = Column(String(), nullable=False)
    password = Column(String(), nullable=False)