import logging

import grpc

from service.gRPC import service_pb2
from service.gRPC import service_pb2_grpc

def main():
  with grpc.insecure_channel("localhost:50051") as channel:
    stub = service_pb2_grpc.UserStub(channel=channel)
    stub.login(service_pb2.LoginRequest(email="test_account", 
      password="not_good_password"))

if __name__ == '__main__':
  main()