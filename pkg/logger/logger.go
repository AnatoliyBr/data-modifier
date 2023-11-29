package logger

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	developmentEncoder = "dev"
	productionEncoder  = "prod"
)

// NewLogger returns zap.Logger with the given
// logFormat ("json", "console"),
// logLevel ("debug", "info", "warn", "error"),
// encoderType ("dev", "prod"),
// outputPaths (e.g. ["stdout"], ["stdout", "./tmp/logs/rpc_traffic.txt"]),
// errorOutputPaths (e.g. ["stderr"], ["stderr", "./tmp/logs/rpc_traffic_errors.txt"]).
func NewLogger(logFormat, logLevel, encoderType string, outputPaths []string, errorOutputPaths []string) (*zap.Logger, error) {
	lvl, err := zapcore.ParseLevel(logLevel)
	if err != nil {
		return nil, err
	}

	rawJSON := []byte(fmt.Sprintf(`{
	"encoding": "%s",
	"level": "%s",
	"encoderConfig": {
		"messageKey": "message",
	    "levelKey": "level",
	    "levelEncoder": "lowercase"
		}
	}`, logFormat, lvl))

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

	cfg.OutputPaths = outputPaths
	cfg.ErrorOutputPaths = errorOutputPaths

	return zap.Must(cfg.Build()), nil
}
