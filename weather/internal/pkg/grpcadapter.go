package grpcadapter

import (
	"log"
	"os"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MaxRetry        = 3
	BackoffDuration = 100 * time.Millisecond
)

// ClientSvcFunc is a function that, given a *grpc.ClientConn, performs
// a gRPC call and returns the response (as any) and an error.
type ClientSvcFunc func(*grpc.ClientConn) (any, error)

// GRPCClientAdapter holds the connection and a standard-library logger.
type GRPCClientAdapter struct {
	conn   *grpc.ClientConn
	logger *log.Logger
}

// New dials the given address with retry middleware and returns an adapter.
// You can pass in your own logger, or nil to use the default logger.
func New(addr string, logger *log.Logger) (*GRPCClientAdapter, error) {
	if logger == nil {
		// Default to standard logger writing to stderr, with date/time prefix
		logger = log.New(os.Stderr, "", log.LstdFlags)
	}

	retryOpts := []retry.CallOption{
		retry.WithMax(MaxRetry),
		retry.WithBackoff(retry.BackoffExponential(BackoffDuration)),
		retry.WithCodes(codes.Unavailable),
	}

	opts := []grpc.DialOption{
		grpc.WithInsecure(), // consider using credentials in production
		grpc.WithChainUnaryInterceptor(
			retry.UnaryClientInterceptor(retryOpts...),
		),
	}

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		logger.Printf("ERROR: failed to dial %s: %v", addr, err)
		return nil, err
	}

	return &GRPCClientAdapter{
		conn:   conn,
		logger: logger,
	}, nil
}

// Do invokes the given ClientSvcFunc with retries; logs success or failure.
func (r *GRPCClientAdapter) Do(f ClientSvcFunc) (any, error) {
	res, err := f(r.conn)
	if s, ok := status.FromError(err); ok {
		// It was a gRPC status error
		switch s.Code() {
		case codes.OK:
			r.logger.Printf(
				"INFO: grpc call success code=%s msg=%q",
				s.Code(), s.Message(),
			)
			return res, nil

		default:
			r.logger.Printf(
				"ERROR: grpc call failed code=%s msg=%q",
				s.Code(), s.Message(),
			)
			return nil, err
		}
	}

	// Not a gRPC status error (maybe a network or client error)
	if err != nil {
		r.logger.Printf(
			"ERROR: unknown error (not gRPC status): %v",
			err,
		)
	}
	return nil, err
}
