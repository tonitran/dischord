package handlers

import (
	"io"
	"log/slog"
	"os"
)

var logger *slog.Logger

func init() {
	if err := os.MkdirAll("logs", 0755); err != nil {
		panic("failed to create logs directory: " + err.Error())
	}
	f, err := os.OpenFile("logs/dischord.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic("failed to open log file: " + err.Error())
	}
	w := io.MultiWriter(os.Stdout, f)
	logger = slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}
