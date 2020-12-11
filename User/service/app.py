from concurrent import futures

from configparser import ConfigParser
import grpc
from tornado.web import Application

from service.db import driver
from service.gRPC import service_pb2_grpc
from service.routes import login, user
from service.server import gRPC

class App():
  def __init__(self, database: driver.Base, logger, config: ConfigParser) -> None:
    self.database = database
    self.logger = logger
    self.config = config
    self.app = self._initalize_server()

  def _initalize_server() -> grpc.Server:
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    service_pb2_grpc.add_UserServicer_to_server(gRPC.Server(), server)
    server.add_insecure_port("[::]:50051")
    return server

  def run(self):
    self.app.start()
    self.app.wait_for_termination()