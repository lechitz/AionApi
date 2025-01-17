package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLoggerSugar initializes a structured logger sending info to stdout and error to stderr
func InitLoggerSugar() (*zap.SugaredLogger, func()) {

	// Encoder JSON padrão do zap
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	// Cria duas LevelEnabler: info e error
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		// Quero todos os logs >= Info e < Error (Info e Warn, por ex.)
		return lvl >= zapcore.InfoLevel && lvl < zapcore.ErrorLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		// Quero todos os logs >= Error
		return lvl >= zapcore.ErrorLevel
	})

	// Destinos
	infoWriter := zapcore.AddSync(os.Stdout)
	errorWriter := zapcore.AddSync(os.Stderr)

	// Cria 2 "cores" distintos
	infoCore := zapcore.NewCore(encoder, infoWriter, infoLevel)
	errorCore := zapcore.NewCore(encoder, errorWriter, errorLevel)

	// Junta num Tee
	core := zapcore.NewTee(infoCore, errorCore)
	logger := zap.New(core).Sugar()

	// Função de cleanup
	cleanup := func() {
		_ = logger.Sync()
	}

	return logger, cleanup
}
