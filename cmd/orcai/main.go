package main

import (
	"OrcAI/internal/core"
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

	defer nc.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	lis, _ := net.Listen("tcp", ":50051")
	grpcServer := grpc.NewServer()

	fallback := "chat"
	bridge := &grpcx.BridgeServer{
		NC: nc,
		Router: core.NewKeywordRouter(map[string]string{
			"how are you": "chat",
			"who are you": "chat",
			"hello":       "chat",
		}, &fallback),
	}

	pb.RegisterBridgeServiceServer(grpcServer, bridge)

	if err = bridge.StartSubscribers(); err != nil {
		logger.Fatal().Err(err).Msg("Failed to start subscribers")
	}

	if err = grpcServer.Serve(lis); err != nil {
		logger.Fatal().Err(err).Msg("Failed to serve gRPC")
	}

	<-c

	if err = lis.Close(); err != nil {
		logger.Fatal().Err(err).Msg("Failed to close listener")
	}
}
