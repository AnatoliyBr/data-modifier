package configs

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		TCP
		WebAPI
		Logger
		App
	}

	App struct {
		MaxClients     int    `toml:"max_clients"`
		NumPoolWorkers uint32 `toml:"num_pool_workers"`
	}

	TCP struct {
		IP   string `toml:"tcp_ip"`
		Port string `toml:"tcp_port"`
	}

	WebAPI struct {
		IP           string `toml:"web_api_ip"`
		Port         string `toml:"web_api_port"`
		Login        string `toml:"login"`
		Password     string `toml:"password"`
		ProtocolType string `toml:"protocol_type"`
		EmployeePath string `toml:"employee_path"`
		AbsencePath  string `toml:"absence_path"`
	}

	Logger struct {
		Format      string `toml:"log_format"`
		Level       string `toml:"log_level"`
		EncoderType string `toml:"encoder_type"`
	}
)

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
		App: App{
			MaxClients:     10,
			NumPoolWorkers: 5,
		},
		TCP: TCP{
			IP:   "127.0.0.1",
			Port: ":8081",
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
			Level:       "debug",
			Format:      "console",
			EncoderType: "dev",
		},
	}

	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
