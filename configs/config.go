// Package configs contains configuration structs.
package configs

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
)

type (
	// Config contains configs of all app layers.
	Config struct {
		WebAPI
		Logger
		GRPC
	}

	// GRPC contains info for gRPC server configuration.
	GRPC struct {
		IP             string `toml:"tcp_ip"`
		Port           string `toml:"tcp_port"`
		MaxClients     int    `toml:"max_clients"`
		NumPoolWorkers uint32 `toml:"num_pool_workers"`
	}

	// WebAPI contains info for WebAPI configuration.
	WebAPI struct {
		IP           string `toml:"web_api_ip"`
		Port         string `toml:"web_api_port"`
		Login        string `toml:"login"`
		Password     string `toml:"password"`
		ProtocolType string `toml:"protocol_type"`
		EmployeePath string `toml:"employee_path"`
		AbsencePath  string `toml:"absence_path"`
	}

	// Logger contains info for logger configuration.
	Logger struct {
		Format          string   `toml:"log_format"`
		Level           string   `toml:"log_level"`
		EncoderType     string   `toml:"encoder_type"`
		OutputPath      []string `toml:"output_path"`
		ErrorOutputPath []string `toml:"error_output_path"`
	}
)

// NewConfig returns app config.
func NewConfig(configPath string) (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	login := os.Getenv("LOGIN")
	password := os.Getenv("PASSWORD")

	if login == "" || password == "" {
		login = "login"
		password = "password"
	}

	config := &Config{
		GRPC: GRPC{
			IP:             "127.0.0.1",
			Port:           ":8081",
			MaxClients:     10,
			NumPoolWorkers: 5,
		},
		WebAPI: WebAPI{
			IP:           "127.0.0.1",
			Port:         ":8082",
			Login:        login,
			Password:     password,
			ProtocolType: "http",
			EmployeePath: "Portal/springApi/api/employees",
			AbsencePath:  "Portal/springApi/api/absences",
		},
		Logger: Logger{
			Level:           "debug",
			Format:          "console",
			EncoderType:     "dev",
			OutputPath:      []string{"stdout"},
			ErrorOutputPath: []string{"stderr"},
		},
	}

	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
