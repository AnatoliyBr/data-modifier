package logger

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

const (
	developmentEncoder = "dev"
	productionEncoder  = "prod"
)

func NewLogger(logFormat, logLevel, encoderType string) (*zap.Logger, error) {
	rawJSON := []byte(fmt.Sprintf(`{
	"level": "%s",
	"encoding": "%s",
	"outputPaths": ["stdout"],
	"errorOutputPaths": ["stderr"],
	"encoderConfig": {
		"messageKey": "message",
	    "levelKey": "level",
	    "levelEncoder": "lowercase"
		}
	}`, logLevel, logFormat))

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}

	switch encoderType {
	case developmentEncoder:
		cfg.EncoderConfig = zap.NewDevelopmentEncoderConfig()
	case productionEncoder:
		cfg.EncoderConfig = zap.NewProductionEncoderConfig()
	}

	return zap.Must(cfg.Build()), nil
}
