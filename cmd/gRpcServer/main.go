package main

import (
	"github.com/Geniuskaa/Bank-system/cmd/gRpcServer/app"
	phoneBalanceV1Pb "github.com/Geniuskaa/Bank-system/pkg/gen/proto/v1"
	"google.golang.org/grpc"
	"net"
	"os"
)

const defaultPort = "9999"
const defaultHost = "0.0.0.0"

func main() {
	port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("APP_HOST")
	if !ok {
		host = defaultHost
	}

	if err := execute(net.JoinHostPort(host, port)); err != nil {
		os.Exit(1)
	}
}

func execute(addr string) (err error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	server := app.NewServer()
	phoneBalanceV1Pb.RegisterPhoneBalancePatternServiceServer(grpcServer, server)

	return grpcServer.Serve(listener)
}
