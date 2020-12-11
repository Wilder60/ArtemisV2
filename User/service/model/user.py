from typing import Any

from sqlalchemy import Column, ForeignKey, Integer, String
from sqlalchemy.orm import relationship

from service.generated import service_pb2
from service.model import schema
from service.model import user

class Info(schema.Base):
    __tablename__ = "info"

    id = Column(Integer, primary_key=True)
    firstname = Column(String(), nullable=False)
    lastname = Column(String(), nullable=False)
    birthday = Column(String(), nullable=True)
    language = Column(String(), nullable=False)
    user_id = Column(Integer, ForeignKey('user.id'))
    user = relationship(user.User)

    @classmethod
    def from_rpc(cls, rpc: service_pb2.UserInfo) -> 'Info':
      return cls(rpc.email, rpc.password, rpc.firstname, rpc.lastname,
                rpc.birthday, rpc.location, rpc.language)
    
    @classmethod
    def from_payload(cls, data: dict[str, Any]) -> 'Info':
      return cls(data["email"], data["password"], data["firstname"], 
                data["lastname"], data["birthday"], data["language"]) 

    def __init__(self, email: str, password: str, firstname: str, 
                lastname: str, dob: str, timezone: str, language: str):
        self.firstname = firstname
        self.lastname = lastname
        self.birthday = dob
        self.language = language
        self.user = user.User(email = email, password = password)

    def to_dict(self) -> dict:
        return self.__dict__