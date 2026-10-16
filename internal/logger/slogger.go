package slogger

import (
	"log/slog"
	"os"
)

type Server struct {
	logger *slog.Logger
}

func New() *slog.Logger {

	handlerOpts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	logger := slog.NewJSONHandler(os.Stderr, handlerOpts)

	return slog.New(logger)
}
