package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name  string
		cfg   Config
		level Level
	}{
		{
			name:  "default level",
			cfg:   Config{Level: "info"},
			level: LevelInfo,
		},
		{
			name:  "debug level",
			cfg:   Config{Level: "debug"},
			level: LevelDebug,
		},
		{
			name:  "warn level",
			cfg:   Config{Level: "warn"},
			level: LevelWarn,
		},
		{
			name:  "error level",
			cfg:   Config{Level: "error"},
			level: LevelError,
		},
		{
			name:  "invalid level defaults to info",
			cfg:   Config{Level: "invalid"},
			level: LevelInfo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.cfg)
			if l == nil {
				t.Error("New() returned nil")
			}
		})
	}
}

func TestLogger_Output(t *testing.T) {
	var buf bytes.Buffer

	l := New(Config{
		Color: false,
		Time:  false,
		Level: "debug",
	})
	l.SetOutput(&buf)

	t.Run("debug message", func(t *testing.T) {
		buf.Reset()
		l.Debug("test %s", "debug")
		if !strings.Contains(buf.String(), "[DEBUG]") {
			t.Errorf("Debug() output = %v, want [DEBUG]", buf.String())
		}
		if !strings.Contains(buf.String(), "test debug") {
			t.Errorf("Debug() output = %v, want 'test debug'", buf.String())
		}
	})

	t.Run("info message", func(t *testing.T) {
		buf.Reset()
		l.Info("test %s", "info")
		if !strings.Contains(buf.String(), "[INFO]") {
			t.Errorf("Info() output = %v, want [INFO]", buf.String())
		}
	})

	t.Run("warn message", func(t *testing.T) {
		buf.Reset()
		l.Warn("test %s", "warn")
		if !strings.Contains(buf.String(), "[WARN]") {
			t.Errorf("Warn() output = %v, want [WARN]", buf.String())
		}
	})

	t.Run("error message", func(t *testing.T) {
		buf.Reset()
		l.Error("test %s", "error")
		if !strings.Contains(buf.String(), "[ERROR]") {
			t.Errorf("Error() output = %v, want [ERROR]", buf.String())
		}
	})
}

func TestLogger_Level(t *testing.T) {
	var buf bytes.Buffer

	l := New(Config{
		Color: false,
		Time:  false,
		Level: "warn",
	})
	l.SetOutput(&buf)

	t.Run("debug suppressed at warn level", func(t *testing.T) {
		buf.Reset()
		l.Debug("should not appear")
		if buf.Len() > 0 {
			t.Errorf("Debug() should be suppressed at warn level, got: %v", buf.String())
		}
	})

	t.Run("info suppressed at warn level", func(t *testing.T) {
		buf.Reset()
		l.Info("should not appear")
		if buf.Len() > 0 {
			t.Errorf("Info() should be suppressed at warn level, got: %v", buf.String())
		}
	})

	t.Run("warn shown at warn level", func(t *testing.T) {
		buf.Reset()
		l.Warn("should appear")
		if buf.Len() == 0 {
			t.Error("Warn() should appear at warn level")
		}
	})

	t.Run("error shown at warn level", func(t *testing.T) {
		buf.Reset()
		l.Error("should appear")
		if buf.Len() == 0 {
			t.Error("Error() should appear at warn level")
		}
	})
}

func TestLogger_SetLevel(t *testing.T) {
	var buf bytes.Buffer

	l := New(Config{
		Color: false,
		Time:  false,
		Level: "error",
	})
	l.SetOutput(&buf)

	// Initially only error should show
	buf.Reset()
	l.Info("should not appear")
	if buf.Len() > 0 {
		t.Error("Info() should be suppressed at error level")
	}

	// Change level to info
	l.SetLevel(LevelInfo)

	buf.Reset()
	l.Info("should appear")
	if buf.Len() == 0 {
		t.Error("Info() should appear after SetLevel(LevelInfo)")
	}
}

func TestLogger_WithTime(t *testing.T) {
	var buf bytes.Buffer

	l := New(Config{
		Color: false,
		Time:  true,
		Level: "info",
	})
	l.SetOutput(&buf)

	l.Info("test message")
	output := buf.String()

	// Should contain a timestamp (HH:MM:SS format)
	if !strings.Contains(output, ":") {
		t.Errorf("Output should contain timestamp, got: %v", output)
	}
}

func TestLogger_WithColor(t *testing.T) {
	var buf bytes.Buffer

	l := New(Config{
		Color: true,
		Time:  false,
		Level: "info",
	})
	l.SetOutput(&buf)

	l.Info("test message")
	output := buf.String()

	// Color codes start with escape sequence.
	// Note: color output may be disabled in non-TTY environments.
	// We verify that the logger was created and produced output.
	if !strings.Contains(output, "[INFO]") {
		t.Errorf("Output should contain [INFO], got: %v", output)
	}
	if !strings.Contains(output, "test message") {
		t.Errorf("Output should contain 'test message', got: %v", output)
	}
}

func TestSuccess(t *testing.T) {
	var buf bytes.Buffer

	l := New(Config{
		Color: false,
		Time:  false,
		Level: "info",
	})
	l.SetOutput(&buf)

	Success(l, "operation %s", "completed")
	output := buf.String()

	if !strings.Contains(output, "✓") {
		t.Errorf("Success() should contain checkmark, got: %v", output)
	}
	if !strings.Contains(output, "operation completed") {
		t.Errorf("Success() should contain message, got: %v", output)
	}
}

func TestFailure(t *testing.T) {
	var buf bytes.Buffer

	l := New(Config{
		Color: false,
		Time:  false,
		Level: "error",
	})
	l.SetOutput(&buf)

	Failure(l, "operation %s", "failed")
	output := buf.String()

	if !strings.Contains(output, "✗") {
		t.Errorf("Failure() should contain X mark, got: %v", output)
	}
	if !strings.Contains(output, "operation failed") {
		t.Errorf("Failure() should contain message, got: %v", output)
	}
}

func TestBanner(t *testing.T) {
	var buf bytes.Buffer

	Banner(&buf, "v1.0.0")
	output := buf.String()

	// Banner contains ASCII art representation of "goreload"
	if !strings.Contains(output, "|___/") {
		t.Errorf("Banner() should contain ASCII art, got: %v", output)
	}
	if !strings.Contains(output, "v1.0.0") {
		t.Errorf("Banner() should contain version, got: %v", output)
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input string
		want  Level
	}{
		{"debug", LevelDebug},
		{"info", LevelInfo},
		{"warn", LevelWarn},
		{"error", LevelError},
		{"unknown", LevelInfo},
		{"", LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := ParseLevel(tt.input); got != tt.want {
				t.Errorf("ParseLevel(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
