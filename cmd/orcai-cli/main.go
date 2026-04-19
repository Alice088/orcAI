package main

import (
	"OrcAI/internal/grpc/pb"
	"context"
	"os"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func main() {
	logger := zerolog.New(os.Stdin).With().Timestamp().Logger()

	root := &cobra.Command{
		Use: "orcai",
	}

	ask := &cobra.Command{
		Use: "ask",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				logger.Fatal().Msg("Empty ask")
			}

			conn, err := grpc.NewClient("localhost:50051", grpc.WithInsecure())
			if err != nil {
				logger.Fatal().Err(err).Msg("Failed to dial orcai")
			}
			defer conn.Close()
			client := pb.NewBridgeServiceClient(conn)
			resp, err := client.Ask(context.Background(), &pb.AskRequest{
				Input: args[0],
			})

			if err != nil {
				logger.Fatal().Err(err).Msg("Failed to request ask")
			}

			logger.Info().Msg(resp.TaskId)
		},
	}

	logs := &cobra.Command{
		Use: "logs",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				logger.Fatal().Msg("Empty logs")
			}

			taskId, err := uuid.Parse(args[0])
			if err != nil {
				logger.Fatal().Err(err).Msg("Invalid task_id")
			}

			conn, err := grpc.NewClient("localhost:50051", grpc.WithInsecure())
			if err != nil {
				logger.Fatal().Err(err).Msg("Failed to dial orcai")
			}
			defer conn.Close()
			client := pb.NewBridgeServiceClient(conn)
			resp, err := client.Logs(context.Background(), &pb.LogRequest{
				TaskId: taskId.String(),
			})

			if err != nil {
				logger.Fatal().Err(err).Msg("Failed to request ask")
			}

			for {
				msg, err := resp.Recv()
				if err != nil {
					break
				}

				logger.Info().Msg(msg.Log)
			}

		},
	}

	root.AddCommand(ask, logs)
	if err := root.Execute(); err != nil {
		logger.Fatal().Err(err).Msg("Failed execute command")
	}
}
