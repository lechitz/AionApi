package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLoggerSugar() (*zap.SugaredLogger, func()) {

	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel && lvl < zapcore.ErrorLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	infoWriter := zapcore.AddSync(os.Stdout)
	errorWriter := zapcore.AddSync(os.Stderr)

	infoCore := zapcore.NewCore(encoder, infoWriter, infoLevel)
	errorCore := zapcore.NewCore(encoder, errorWriter, errorLevel)

	core := zapcore.NewTee(infoCore, errorCore)
	logger := zap.New(core).Sugar()

	cleanup := func() {
		_ = logger.Sync()
	}

	return logger, cleanup
}
