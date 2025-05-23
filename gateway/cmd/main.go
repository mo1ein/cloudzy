package main

import (
	"context"
	"fmt"
	"gateway/cmd/command"
	"gateway/internal/config"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"log"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	const description = "ServiceName"
	root := &cobra.Command{Short: description}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	//app.InitLogger()

	location, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		log.Fatal(err)
	}
	time.Local = location

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var serverCmd command.Server
	//var grpc command.GRPCServer
	root.AddCommand(
		// todo: add logger here
		serverCmd.Command(ctx, zerolog.Logger{}, cfg),
		//mysqlCommand.Migrate{}.Command(ctx, &cfg.Database.MySQL),
		//grpc.Command(ctx, loggerGRPC, cfg),
	)

	if err := root.Execute(); err != nil {
		// todo: correct log.fatal
		log.Fatal(fmt.Sprintf("failed to execute root command: \n%v", err))
	}
}
