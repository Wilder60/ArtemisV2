"""Data class for SQLAlchemy to create login table"""
from sqlalchemy import Column, Integer, String
from typing import Any

from service.generated import service_pb2
from service.model import schema

class Request(schema.Base):
    __tablename__ = "login"

    id = Column(Integer, primary_key=True)
    email = Column(String(), nullable=False)
    password = Column(String(), nullable=False)

    @classmethod
    def from_rpc(cls, rpc: service_pb2.LoginRequest) -> 'Request':
        return cls(rpc.email, rpc.password)

    @classmethod
    def from_payload(cls, data: dict[str, Any]) -> 'Request':
        return cls(data["email"], data["password"])

    def __init__(self, email: str, password: str) -> 'Request':
        self.email = email
        self.password = password

    def __str__(self) -> str:
        return self.__repr__()

    def __repr__(self) -> str:
        return "id:{}\temail:{}\tpassword:{}".format(
            self.id, self.email, self.password)
