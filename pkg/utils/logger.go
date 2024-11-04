package utils

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger *zap.Logger
	once   sync.Once
)

type LogConfig struct {
	Level      string
	OutputPath string
	MaxSize    int  // megabytes
	MaxBackups int  // number of backups
	MaxAge     int  // days
	Compress   bool // compress old files
}

func DefaultLogConfig() *LogConfig {
	return &LogConfig{
		Level:      "info",
		OutputPath: "logs/gateway.log",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}
}

// InitLogger initializes the global logger instance
func InitLogger(config *LogConfig) *zap.Logger {
	once.Do(func() {
		// Create directory if it doesn't exist
		os.MkdirAll("logs", 0744)

		// Configure logging output
		writer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   config.OutputPath,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
		})

		// Configure encoding
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		// Determine log level
		var level zapcore.Level
		switch config.Level {
		case "debug":
			level = zapcore.DebugLevel
		case "info":
			level = zapcore.InfoLevel
		case "warn":
			level = zapcore.WarnLevel
		case "error":
			level = zapcore.ErrorLevel
		default:
			level = zapcore.InfoLevel
		}

		// Create core
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), writer),
			level,
		)

		// Create logger
		logger = zap.New(core,
			zap.AddCaller(),
			zap.AddCallerSkip(1),
			zap.AddStacktrace(zapcore.ErrorLevel),
		)
	})

	return logger
}

// GetLogger returns the global logger instance
func GetLogger() *zap.Logger {
	if logger == nil {
		return InitLogger(DefaultLogConfig())
	}
	return logger
}

// Log levels
func Info(msg string, fields ...zapcore.Field) {
	GetLogger().Info(msg, fields...)
}

func Debug(msg string, fields ...zapcore.Field) {
	GetLogger().Debug(msg, fields...)
}

func Warn(msg string, fields ...zapcore.Field) {
	GetLogger().Warn(msg, fields...)
}

func Error(msg string, fields ...zapcore.Field) {
	GetLogger().Error(msg, fields...)
}

func Fatal(msg string, fields ...zapcore.Field) {
	GetLogger().Fatal(msg, fields...)
}

// Custom field helpers
func Fields(keysAndValues ...interface{}) []zapcore.Field {
	fields := make([]zapcore.Field, 0, len(keysAndValues)/2)
	for i := 0; i < len(keysAndValues); i += 2 {
		key, ok := keysAndValues[i].(string)
		if !ok || i+1 >= len(keysAndValues) {
			continue
		}
		fields = append(fields, zap.Any(key, keysAndValues[i+1]))
	}
	return fields
}
