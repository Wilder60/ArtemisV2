import bcrypt
import grpc

from service.db import driver
from service.generated import service_pb2
from service.generated import service_pb2_grpc
from service.model import login
from service.model import user

class Server(service_pb2_grpc.UserServicer):
  def __init__(self, database: driver.Base) -> None:
    self.database = database

  def create_user(self, request: service_pb2.UserInfo, 
                  context: grpc.ServicerContext) -> service_pb2.Empty:
    """Creates a user account

    Args:
      request: The request data in the form of a UserInfo instance
      context: The rpc request context
    
    Returns:
      The Empty object
    """
    request = user.Info.from_rpc(request)
    self.database.insert_user(request)
    return service_pb2.Empty()

  def get_user_info(self, request: service_pb2.GetUser, 
                    context:grpc.ServicerContext) -> service_pb2.UserInfo:
    return

  def modify_user(self, request: service_pb2.Empty, 
                  context:grpc.ServicerContext) -> service_pb2.Empty:
    return

  def delete_user(self, request: service_pb2.Empty, 
                  context: grpc.ServicerContext) -> service_pb2.Empty:
    return

  def login(self, request: service_pb2.LoginRequest, 
                  context: grpc.ServicerContext) -> service_pb2.Empty:
    print(request.email)
    print(request.password)

    if not bcrypt.checkpw(request.password, ""):
      context.set_code(grpc.StatusCode.PERMISSION_DENIED)
      context.set_details("Invalid password")
      return service_pb2.Empty()

    return

  def update_password(self, request: service_pb2.PasswordUpdate, 
                      context: grpc.ServicerContext) -> service_pb2.Empty:
    return