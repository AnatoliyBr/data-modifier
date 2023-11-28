// Package grpcserver implements gRPC server.
package grpcserver

import (
	"context"
	"fmt"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCServer is a custom gRPC server with RPC requests logger.
type GRPCServer struct {
	server *grpc.Server
	notify chan error
}

// NewGRPCServer returns GRPCServer with the given number
// of goroutines for traffic processing.
func NewGRPCServer(NumPoolWorkers uint32) *GRPCServer {
	numPoolWorkersOpt := grpc.NumStreamWorkers(NumPoolWorkers)

	// loggingOpts configure the interseptor for logging
	// RPC requests.
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived, logging.PayloadSent,
		),
	}

	// recoveryOpts configure the interseptor for restore and
	// process the panic if it happens inside the handler.
	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			zap.L().Error(fmt.Sprintf("Recovered from panic: %v", p))
			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	return &GRPCServer{
		server: grpc.NewServer(numPoolWorkersOpt, grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(InterceptorLogger(zap.L()), loggingOpts...),
			recovery.UnaryServerInterceptor(recoveryOpts...),
		)),
		notify: make(chan error, 1),
	}
}

// StartGRPCServer calls method for serving connection
// on the given listener and accepts err from it.
func (s *GRPCServer) StartGRPCServer(l net.Listener) {
	go func() {
		s.notify <- s.server.Serve(l)
		close(s.notify)
	}()
}

// Notify returns notify channel field.
func (s *GRPCServer) Notify() <-chan error {
	return s.notify
}

// GracefulShutdown stops the gRPC server gracefully.
func (s *GRPCServer) GracefulShutdown() {
	s.server.GracefulStop()
}

// GetServer returns gRPC server field.
func (s *GRPCServer) GetServer() *grpc.Server {
	return s.server
}

// InterceptorLogger adapts zap.Logger to logging.Logger.
func InterceptorLogger(l *zap.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		convertLevel := zapcore.Level(lvl / 4)

		metadata := make(map[string]string, (len(fields)-2)/2)
		for i := 0; i < len(fields)-3; i = i + 2 {
			metadata[fields[i].(string)] = fields[i+1].(string)
		}

		traffic := fields[len(fields)-1]

		convertMetadata := zap.Any("metadata", metadata)
		convertTraffic := zap.Any("traffic", traffic)

		l.Log(convertLevel, msg, convertMetadata, convertTraffic)
	})
}
