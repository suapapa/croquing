package logging

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

func TestNewLoggerTextFormat(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := NewLogger(FormatText, &buf)
	logger.Info("hello", "key", "value")

	if !strings.Contains(buf.String(), "hello") {
		t.Fatalf("log output = %q, want message", buf.String())
	}
	if strings.Contains(buf.String(), `"key"`) {
		t.Fatalf("log output = %q, want text format not json", buf.String())
	}
}

func TestNewLoggerJSONFormat(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := NewLogger(FormatJSON, &buf)
	logger.Info("hello", "key", "value")

	if !strings.Contains(buf.String(), `"msg":"hello"`) {
		t.Fatalf("log output = %q, want json message field", buf.String())
	}
}

func TestInitSetsDefaultLogger(t *testing.T) {
	var buf bytes.Buffer
	slog.SetDefault(NewLogger(FormatText, &buf))
	slog.Info("default logger")

	if !strings.Contains(buf.String(), "default logger") {
		t.Fatalf("log output = %q, want message", buf.String())
	}
}
