package logger

import (
	"log/slog"
	"os"
)

var (
	log   *slog.Logger
	level string = "info"
)

func Init() {
	lvl := getSlogLevel(level)
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	})
	log = slog.New(handler)

}

func getSlogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func Debug(msg string, args ...interface{}) {
	log.Debug(msg, args...)
}

func Info(msg string, args ...interface{}) {
	log.Info(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	log.Warn(msg, args...)
}

func Error(msg string, args ...interface{}) {
	log.Error(msg, args...)
}
