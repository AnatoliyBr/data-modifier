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
	cfg, err := configs.NewConfig(configPath)
	if err != nil {
		return err
	}

	// Logger
	log, err := logger.NewLogger(
		cfg.Logger.Format,
		cfg.Logger.Level,
		cfg.Logger.EncoderType)
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
		EmployeeURL: fmt.Sprintf("https://%s%s/%s",
			cfg.WebAPI.IP,
			cfg.WebAPI.Port,
			cfg.WebAPI.EmployeePath),
		AbsenceURL: fmt.Sprintf("https://%s%s/%s",
			cfg.WebAPI.IP,
			cfg.WebAPI.Port,
			cfg.WebAPI.AbsencePath),
	}

	if err = cred.Validate(); err != nil {
		return err
	}

	webAPI := webapi.NewUserWebAPI(cred)

	// UseCase
	zap.L().Info("Initializing AppUseCase...")
	uc := usecase.NewAppUseCase(webAPI)

	// Controller
	zap.L().Info("Initializing NewDataModifierService...")
	src := datamodifier.NewDataModifierService(uc)

	// GRPC server
	zap.L().Info("Initializing GRPCServer...")
	zap.L().Debug(fmt.Sprintf("Listen %s%s", cfg.TCP.IP, cfg.TCP.Port))
	l, err := net.Listen("tcp", cfg.TCP.IP+cfg.TCP.Port)
	if err != nil {
		return err
	}

	zap.L().Debug(fmt.Sprintf("MaxClients: %d", cfg.App.MaxClients))
	l = netutil.LimitListener(l, cfg.App.MaxClients)

	zap.L().Debug(fmt.Sprintf("NumPoolWorkers: %d", cfg.App.NumPoolWorkers))
	s := grpcserver.NewGRPCServer(cfg.App.NumPoolWorkers)

	v1.RegisterDataModifierServer(s.GetServer(), src)

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
