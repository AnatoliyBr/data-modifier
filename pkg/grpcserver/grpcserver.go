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

type GRPCServer struct {
	server *grpc.Server
	notify chan error
}

func NewGRPCServer(NumPoolWorkers uint32) *GRPCServer {
	numPoolWorkersOpt := grpc.NumStreamWorkers(NumPoolWorkers)

	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived, logging.PayloadSent,
		),
	}

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

func (s *GRPCServer) StartGRPCServer(l net.Listener) {
	go func() {
		s.notify <- s.server.Serve(l)
		close(s.notify)
	}()
}

func (s *GRPCServer) Notify() <-chan error {
	return s.notify
}

func (s *GRPCServer) GracefulShutdown() {
	s.server.GracefulStop()
}

func (s *GRPCServer) GetServer() *grpc.Server {
	return s.server
}

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
