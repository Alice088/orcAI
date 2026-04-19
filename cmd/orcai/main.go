package main

import (
	grpcx "OrcAI/internal/grpc"
	"OrcAI/internal/grpc/pb"
	"OrcAI/pkg/cfg"
	"OrcAI/pkg/logging"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
)

func main() {
	logger := logging.MustLogger()
	config := cfg.MustConfig(logger)

	nc, err := nats.Connect(config.NATS.URL)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect NATS")
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	lis, _ := net.Listen("tcp", ":50051")
	grpcServer := grpc.NewServer()

	pb.RegisterBridgeServiceServer(grpcServer, &grpcx.BridgeServer{
		NC: nc,
	})

	if err = grpcServer.Serve(lis); err != nil {
		logger.Fatal().Err(err).Msg("Failed to serve gRPC")
	}

	<-c

	if err = lis.Close(); err != nil {
		logger.Fatal().Err(err).Msg("Failed to close listener")
	}
}
