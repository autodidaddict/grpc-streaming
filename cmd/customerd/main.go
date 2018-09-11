package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/autodidaddict/grpc-streaming/internal/config"
	"github.com/autodidaddict/grpc-streaming/internal/handler"
	"github.com/autodidaddict/grpc-streaming/internal/logging"
	pb "github.com/autodidaddict/grpc-streaming/proto"
	"github.com/go-kit/kit/log/level"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	logger := logging.NewLogger(config.SERVICE_NAME, config.Version)

	rpcHandler := handler.NewCustomersHandler(logger)

	port := flag.Int("port", 8088, "Service Port Number")
	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		level.Error(logger).Log("msg", "TCP Failure", "error", err)
	} else {
		level.Info(logger).Log("msg", "Customer Service listening", "port", *port)
		grpcServer := grpc.NewServer()
		pb.RegisterCustomersServer(grpcServer, rpcHandler)
		reflection.Register(grpcServer)
		grpcServer.Serve(listener)
	}
}
