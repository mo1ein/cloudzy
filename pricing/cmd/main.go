package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"pricing/cmd/command"
	"pricing/internal/config"
)

func main() {
	const description = "ServiceName"
	root := &cobra.Command{Short: description}

	// Initialize zerolog with console output and timestamps
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("failed to load configuration")
	}

	// Set local timezone
	location, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("failed to set timezone")
	}
	time.Local = location

	// Handle OS signals for graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Wire up your gRPC server command, passing the zerolog.Logger
	var grpcCmd command.GRPCServer
	var processorCmd command.Worker
	root.AddCommand(
		grpcCmd.Command(ctx, logger, cfg),
		processorCmd.Command(ctx, logger, cfg),
	)

	// Execute root command
	if err := root.Execute(); err != nil {
		logger.Fatal().
			Err(err).
			Msg("failed to execute root command")
	}
}
