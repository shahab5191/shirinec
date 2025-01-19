package utils

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

func InitLogger() {
    encoderConfig := zapcore.EncoderConfig{
        TimeKey:        "time",
        LevelKey:       "level",
        NameKey:        "logger",
        CallerKey:      "caller",
        MessageKey:     "msg",
        StacktraceKey:  "stacktrace",
        LineEnding:     zapcore.DefaultLineEnding,
        EncodeLevel:    zapcore.CapitalColorLevelEncoder, // Colored output
        EncodeTime:     zapcore.ISO8601TimeEncoder,       // ISO8601 time format
        EncodeDuration: zapcore.StringDurationEncoder,
        EncodeCaller:   zapcore.ShortCallerEncoder,
    }

    consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

    core := zapcore.NewCore(
        consoleEncoder,
        zapcore.AddSync(os.Stdout),
        zapcore.DebugLevel,
    )

	logger := zap.New(core, zap.AddCaller())
	Logger = logger.Sugar()
}
