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
		IP       string `toml:"web_api_ip"`
		Port     string `toml:"web_api_port"`
		Login    string `toml:"login"`
		Password string `toml:"password"`
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
			IP:   "localhost",
			Port: ":8080",
		},
		WebAPI: WebAPI{
			Login:    login,
			Password: password,
		},
	}

	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
