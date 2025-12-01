package logger

import (
	"os"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)	

func CreateZapLogger() *zap.SugaredLogger {
	encoderCfg := zap.NewProductionEncoderConfig()

	if os.Getenv("APP_ENV") == "development" {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	}


	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	config := zap.Config{
		Level: zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		DisableCaller: true,
		DisableStacktrace: true,
		Sampling: nil,
		Encoding: "json",
		EncoderConfig: encoderCfg,
		OutputPaths: []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields: map[string]interface{}{
			"pid": os.Getpid(),
		},

	}

	return zap.Must(config.Build()).Sugar()
}