package log

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"testing"
)

func TestLog(t *testing.T) {
	levels := [4]slog.Level{slog.LevelInfo, slog.LevelDebug, slog.LevelWarn, slog.LevelError}
	for _, level := range levels {
		Log("Message", level)
	}
}

func TestLogInfo(t *testing.T) {
	Info("Info message")
}

func TestLogError(t *testing.T) {
	Error("Error message", status.Error(codes.Internal, "Test error"))
}

func TestLogDebug(t *testing.T) {
	Warning("Warning message")
}
