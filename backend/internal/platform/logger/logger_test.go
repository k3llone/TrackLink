package logger

import (
	"io"
	"os"
	"strings"
	"testing"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}

	os.Stdout = w
	defer func() {
		_ = w.Close()
		os.Stdout = old
	}()

	fn()

	_ = w.Close()
	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read stdout: %v", err)
	}
	return string(out)
}

func TestNewTextFormat(t *testing.T) {
	out := captureStdout(t, func() {
		log := New("text")
		log.Info("format_test")
	})

	if !strings.Contains(out, "level=INFO") {
		t.Fatalf("expected text log level, got %q", out)
	}
	if !strings.Contains(out, "msg=format_test") {
		t.Fatalf("expected text log message, got %q", out)
	}
}

func TestNewJSONFormat(t *testing.T) {
	out := captureStdout(t, func() {
		log := New("json")
		log.Info("format_test")
	})

	if !strings.Contains(out, "\"level\":\"INFO\"") {
		t.Fatalf("expected json log level, got %q", out)
	}
	if !strings.Contains(out, "\"msg\":\"format_test\"") {
		t.Fatalf("expected json log message, got %q", out)
	}
}
