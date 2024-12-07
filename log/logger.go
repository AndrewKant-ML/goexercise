package log

import (
	"fmt"
	"log/slog"
)

// Log logs a message with a given level of
func Log(message string, level slog.Level) {
	switch level {
	case slog.LevelInfo:
		slog.Info(message)
	case slog.LevelDebug:
		slog.Debug(message)
	case slog.LevelWarn:
		slog.Warn(message)
	case slog.LevelError:
		slog.Error(message)
	}
}

func Info(msg string) {
	Log(msg, slog.LevelInfo)
}

func Warning(msg string) {
	Log(msg, slog.LevelWarn)
}

func ErrorMessage(msg string) {
	Log(msg, slog.LevelError)
}

func Error(msg string, err error) {
	ErrorMessage(fmt.Sprintf("%s. Full error: %v", msg, err))
}
