package app

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/AnatoliyBr/data-modifier/configs"
	"github.com/AnatoliyBr/data-modifier/internal/controller/datamodifier"
	v1 "github.com/AnatoliyBr/data-modifier/pkg/api/v1"
	"github.com/AnatoliyBr/data-modifier/pkg/grpcserver"
	"github.com/AnatoliyBr/data-modifier/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/net/netutil"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "configs/config.toml", "path to config file")
}

func Run() error {
	// Config
	flag.Parse()
	config, err := configs.NewConfig(configPath)
	if err != nil {
		return err
	}

	// Logger
	log, err := logger.NewLogger(config.Logger.Format, config.Logger.Level, config.Logger.EncoderType)
	defer log.Sync()

	if err != nil {
		return err
	}

	restore := zap.ReplaceGlobals(log)
	defer restore()

	// WebAPI

	// UseCase

	// Controller
	srv := &datamodifier.DataModifierService{}

	// GRPC server
	zap.L().Debug(fmt.Sprintf("Listen port: %s", config.TCP.Port))
	l, err := net.Listen("tcp", config.TCP.IP+config.TCP.Port)
	if err != nil {
		return err
	}

	l = netutil.LimitListener(l, config.App.MaxClients)

	s := grpcserver.NewGRPCServer(config.App.NumPoolWorkers)

	v1.RegisterDataModifierServer(s.GetServer(), srv)

	s.StartGRPCServer(l)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case signal := <-interrupt:
		zap.L().Info("app - Run - signal: " + signal.String())
	case err = <-s.Notify():
		zap.L().Error("app - Run - GRPCServer.Notify: %w", zap.Error(err))
	}

	// Graceful shutdown
	s.GracefulShutdown()
	zap.L().Info("app - Run - GRPCServer.GracefulShutdown")

	return nil
}
