package logging

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

const (
	FormatText = "text"
	FormatJSON = "json"
)

// Init configures the default slog logger on stdout.
func Init(format string) {
	normalized := strings.ToLower(strings.TrimSpace(format))
	if normalized != FormatText && normalized != FormatJSON && normalized != "" {
		fmt.Fprintf(os.Stderr, "unknown log format %q, using text\n", format)
		normalized = FormatText
	}
	slog.SetDefault(NewLogger(normalized, os.Stdout))
}

// NewLogger builds a slog logger with the given output format.
func NewLogger(format string, w io.Writer) *slog.Logger {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case FormatJSON:
		return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		return slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
}
