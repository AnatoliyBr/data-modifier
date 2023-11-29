// Package app configures and runs application.
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
	"github.com/AnatoliyBr/data-modifier/internal/entity"
	"github.com/AnatoliyBr/data-modifier/internal/usecase"
	"github.com/AnatoliyBr/data-modifier/internal/webapi"
	v1 "github.com/AnatoliyBr/data-modifier/pkg/api/v1"
	"github.com/AnatoliyBr/data-modifier/pkg/grpcserver"
	"github.com/AnatoliyBr/data-modifier/pkg/logger"
	"github.com/AnatoliyBr/data-modifier/pkg/testserver"
	"go.uber.org/zap"
	"golang.org/x/net/netutil"
)

var (
	configPath string
	testServer bool
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/config.toml", "path to config file")
	flag.BoolVar(&testServer, "test-server", false, "flag for starting test server")
}

// Run creates objects via constructors.
func Run() error {
	// Config
	flag.Parse()
	cfg, err := configs.NewConfig(configPath)
	if err != nil {
		return err
	}

	// Logger
	log, err := logger.NewLogger(
		cfg.Logger.Format,
		cfg.Logger.Level,
		cfg.Logger.EncoderType,
		cfg.Logger.OutputPath,
		cfg.Logger.ErrorOutputPath)
	defer log.Sync()

	if err != nil {
		return err
	}

	restore := zap.ReplaceGlobals(log)
	defer restore()

	// WebAPI
	zap.L().Info("Initializing UserWebAPI...")
	cred := &entity.Credentials{
		IP:       cfg.WebAPI.IP,
		Port:     cfg.WebAPI.Port,
		Login:    cfg.WebAPI.Login,
		Password: cfg.WebAPI.Password,
		EmployeeURL: fmt.Sprintf("%s://%s%s/%s",
			cfg.WebAPI.ProtocolType,
			cfg.WebAPI.IP,
			cfg.WebAPI.Port,
			cfg.WebAPI.EmployeePath),
		AbsenceURL: fmt.Sprintf("%s://%s%s/%s",
			cfg.WebAPI.ProtocolType,
			cfg.WebAPI.IP,
			cfg.WebAPI.Port,
			cfg.WebAPI.AbsencePath),
	}

	if err = cred.Validate(); err != nil {
		return err
	}

	webAPI := webapi.NewUserWebAPI(cred)

	// TestServer
	if testServer {
		ts := testserver.NewTestServer(
			cfg.WebAPI.EmployeePath,
			cfg.WebAPI.AbsencePath,
			cfg.WebAPI.Port,
			webAPI.BasicAuthToken,
			entity.TestUser(),
			entity.TestUserAbsenceData(),
		)

		zap.L().Debug(fmt.Sprintf("Starting TestServer on port %s", cfg.WebAPI.Port))
		go func() error {
			return ts.StartTestServer()
		}()
	}

	// UseCase
	zap.L().Info("Initializing AppUseCase...")
	uc := usecase.NewAppUseCase(webAPI)

	// Controller
	zap.L().Info("Initializing NewDataModifierService...")
	src := datamodifier.NewDataModifierService(uc)

	// GRPC server
	zap.L().Info("Initializing GRPCServer...")
	l, err := net.Listen("tcp", cfg.GRPC.IP+cfg.GRPC.Port)
	if err != nil {
		return err
	}

	zap.L().Debug(fmt.Sprintf("MaxClients: %d", cfg.GRPC.MaxClients))
	l = netutil.LimitListener(l, cfg.GRPC.MaxClients)

	zap.L().Debug(fmt.Sprintf("NumPoolWorkers: %d", cfg.GRPC.NumPoolWorkers))
	s := grpcserver.NewGRPCServer(cfg.GRPC.NumPoolWorkers)

	v1.RegisterDataModifierServer(s.GetServer(), src)

	zap.L().Debug(fmt.Sprintf("Starting GRPCServer on port %s", cfg.GRPC.Port))
	s.StartGRPCServer(l)

	// Waiting signal
	zap.L().Info("Configuring graceful shutdown...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case signal := <-interrupt:
		zap.L().Info("app - Run - signal: " + signal.String())
	case err = <-s.Notify():
		zap.L().Error("app - Run - GRPCServer.Notify: %w", zap.Error(err))
	}

	// Graceful shutdown
	zap.L().Info("Shutting down...")
	s.GracefulShutdown()

	return nil
}
