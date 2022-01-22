// Package log2 knows how to log
package log2

import (
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

// New returns a new zap.SugaredLogger
func New() (*zap.SugaredLogger, error) {
	logType, ok := os.LookupEnv("LOG_TYPE")
	if !ok {
		logType = "JSON"
	}

	var level zapcore.Level

	levelString, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		level = zap.InfoLevel
	} else if strings.ToLower(levelString) == "debug" {
		level = zap.DebugLevel
		fmt.Println("DEBUG_LEVEL: debug")
	}

	switch strings.ToLower(logType) {
	case "JSON":
		l, err := zap.NewProduction()
		if err != nil {
			return nil, fmt.Errorf("could not get prod logger: %w", err)
		}

		return l.Sugar(), nil
	default:
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeTime = func(time.Time, zapcore.PrimitiveArrayEncoder) {}
		config.EncoderConfig.EncodeCaller = func(zapcore.EntryCaller, zapcore.PrimitiveArrayEncoder) {}
		config.Level.SetLevel(level)

		l, err := config.Build()
		if err != nil {
			return nil, fmt.Errorf("creating logger development config: %w", err)
		}

		return l.Sugar(), nil
	}
}
