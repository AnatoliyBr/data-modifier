package grpcserver

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	server *grpc.Server
	notify chan error
}

func NewGRPCServer(NumPoolWorkers uint32) *GRPCServer {
	numPoolWorkersOpt := grpc.NumStreamWorkers(NumPoolWorkers)
	return &GRPCServer{
		server: grpc.NewServer(numPoolWorkersOpt),
		notify: make(chan error, 1),
	}
}

func (s *GRPCServer) StartGRPCServer(l net.Listener) {
	log.Print("grpcserver - StartGRPCServer")
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
