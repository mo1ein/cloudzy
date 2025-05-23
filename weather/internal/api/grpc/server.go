package grpcapi

import (
	"context"
	pb "github.com/mo1ein/cloudzy/proto/gen"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"net"
	"time"
	grpchandlers "weather/internal/api/grpc/handler"
)

const (
	MaxConnectionAge      = time.Second * 30
	MaxConnectionAgeGrace = time.Second * 10
)

type Server struct {
	listener net.Listener
	pb.UnimplementedWeatherServiceServer
	grpcServer *grpc.Server
	handler    grpchandlers.Handler
	logger     zerolog.Logger
}

func New(listener net.Listener, handler grpchandlers.Handler, logger zerolog.Logger) Server {
	// init interceptors
	unaryServerInterceptors := NewUnaryServerInterceptor(logger, GRPCLoggerConfig{
		MoreLogStatusStart: int(codes.InvalidArgument),
	})

	var opts []grpc.ServerOption
	opts = append(
		opts,
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAge:      MaxConnectionAge,
			MaxConnectionAgeGrace: MaxConnectionAgeGrace,
		}),
		grpc.ChainUnaryInterceptor(
			unaryServerInterceptors.RecoverInterceptor(),
			unaryServerInterceptors.UnaryLogInterceptor(),
		),
	)
	grpcSrv := grpc.NewServer(opts...)

	srv := Server{
		listener:                          listener,
		UnimplementedWeatherServiceServer: pb.UnimplementedWeatherServiceServer{},
		grpcServer:                        grpcSrv,
		logger:                            logger,
		handler:                           handler,
	}

	pb.RegisterWeatherServiceServer(grpcSrv, &srv)

	return srv
}
func (r *Server) Start(ctx context.Context) error {
	// Corrected logging syntax
	r.logger.Info().
		Str("address", r.listener.Addr().String()).
		Msg("grpc server starting")

	srvError := make(chan error)
	go func() {
		srvError <- r.grpcServer.Serve(r.listener)
	}()

	select {
	case <-ctx.Done():
		// Corrected shutdown logging
		r.logger.Info().Msg("grpc server is shutting down")
		r.grpcServer.GracefulStop()
	case err := <-srvError:
		return err
	}

	return nil
}
