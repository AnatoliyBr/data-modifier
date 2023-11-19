package app

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/AnatoliyBr/data-modifier/configs"
	v1 "github.com/AnatoliyBr/data-modifier/pkg/api/v1"
	"github.com/AnatoliyBr/data-modifier/pkg/datamodifier"
	"github.com/AnatoliyBr/data-modifier/pkg/grpcserver"
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

	// WebAPI

	// UseCase ?

	// Controller
	l, err := net.Listen("tcp", config.TCP.IP+config.TCP.Port)
	if err != nil {
		return err
	}

	l = netutil.LimitListener(l, config.App.MaxClients)
	s := grpcserver.NewGRPCServer(config.App.NumPoolWorkers)

	srv := &datamodifier.DataModifierService{WebAPI: struct {
		IP       string
		Port     string
		Login    string
		Password string
	}{
		IP:       config.WebAPI.IP,
		Port:     config.WebAPI.Port,
		Login:    config.WebAPI.Login,
		Password: config.WebAPI.Password,
	}}

	v1.RegisterDataModifierServer(s.GetServer(), srv)

	s.StartGRPCServer(l)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case signal := <-interrupt:
		log.Print("app - Run - signal: " + signal.String())
	case err = <-s.Notify():
		log.Print("app - Run - GRPCServer.Notify: %w", err)
	}

	// Shutdown
	s.GracefulShutdown()
	log.Print("app - Run - GRPCServer.GracefulShutdown")

	return nil
}
