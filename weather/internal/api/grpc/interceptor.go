package grpcapi

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCLoggerConfig struct {
	MoreLogStatusStart int
}

type UnaryServerInterceptor struct {
	logger           zerolog.Logger
	GRPCLoggerConfig GRPCLoggerConfig
}

func NewUnaryServerInterceptor(logger zerolog.Logger, grpcLoggerConfig GRPCLoggerConfig) *UnaryServerInterceptor {
	return &UnaryServerInterceptor{
		logger:           logger,
		GRPCLoggerConfig: grpcLoggerConfig,
	}
}

func (r *UnaryServerInterceptor) UnaryLogInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		latency := time.Since(start)

		st := status.Convert(err)
		code := st.Code()

		logEvent := r.logger.Info()
		if code >= codes.Internal {
			logEvent = r.logger.Error()
		} else if err != nil {
			logEvent = r.logger.Error()
		}

		logEvent.
			Str("method", info.FullMethod).
			Int("status_code", int(code)).
			Str("latency", latency.String()).
			Str("error", st.Message()).
			Msg("gRPC request")

		return resp, err
	}
}

func (r *UnaryServerInterceptor) RecoverInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		defer func() {
			if recovered := recover(); recovered != nil {
				r.logger.Error().
					Interface("error", recovered).
					Str("stack", string(debug.Stack())).
					Msg("PANIC recovered")
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()
		return handler(ctx, req)
	}
}
