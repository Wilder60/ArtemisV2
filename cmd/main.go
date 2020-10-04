package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Wilder60/KeyRing/adpater"
	"github.com/Wilder60/KeyRing/internal/grpc/service"

	"google.golang.org/grpc"
)

func getNetListener(port uint) net.Listener {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal("failed to listen:" + err.Error())
		panic(fmt.Sprintf("failed to listen: %s", err.Error()))
	}

	return lis
}

func main() {
	netListener := getNetListener(7000)
	gRPCServer := grpc.NewServer()

	serviceImpl := adpater.NewServiceGrpcImpl()
	service.RegisterKeyHookServiceServer(gRPCServer, serviceImpl)

}
