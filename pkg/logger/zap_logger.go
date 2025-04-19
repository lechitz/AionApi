package logger

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	failedToFlushLogger = "Failed to flush logger: %v"
)

func NewZapLogger() (*zap.SugaredLogger, func()) {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	encoder := zapcore.NewJSONEncoder(encoderCfg)

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel && lvl < zapcore.ErrorLevel
	})

	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	infoWriter := zapcore.Lock(os.Stdout)
	errorWriter := zapcore.Lock(os.Stderr)

	infoCore := zapcore.NewCore(encoder, infoWriter, infoLevel)
	errorCore := zapcore.NewCore(encoder, errorWriter, errorLevel)

	tee := zapcore.NewTee(infoCore, errorCore)

	logger := zap.New(tee, zap.AddCaller(), zap.AddCallerSkip(1))
	sugar := logger.Sugar()

	cleanup := func() {
		if err := sugar.Sync(); err != nil {
			log.Printf(failedToFlushLogger, err)
		}
	}

	return sugar, cleanup
}
