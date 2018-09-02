package main

import (
	"flag"
	"fmt"
	"github.com/autodidaddict/grpc-streaming/internal/config"
	"github.com/autodidaddict/grpc-streaming/internal/handler"
	"github.com/autodidaddict/grpc-streaming/internal/logging"
	"github.com/go-kit/kit/log/level"
	"google.golang.org/grpc"
	"net"
	pb "github.com/autodidaddict/grpc-streaming/proto"
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
		grpcServer.Serve(listener)
	}
}