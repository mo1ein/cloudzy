package grpcapi

import (
	"context"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"net"
	grpchandlers "pricing/internal/api/grpc/handler"
	pb "pricing/proto/gen"
	"time"
)

const (
	MaxConnectionAge      = time.Second * 30
	MaxConnectionAgeGrace = time.Second * 10
)

type Server struct {
	listener net.Listener
	// todo: changed
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
func (s *Server) Start(ctx context.Context) error {
	// Corrected logging syntax
	s.logger.Info().
		Str("address", s.listener.Addr().String()).
		Msg("grpc server starting")

	srvError := make(chan error)
	go func() {
		srvError <- s.grpcServer.Serve(s.listener)
	}()

	select {
	case <-ctx.Done():
		// Corrected shutdown logging
		s.logger.Info().Msg("grpc server is shutting down")
		s.grpcServer.GracefulStop()
	case err := <-srvError:
		return err
	}

	return nil
}
