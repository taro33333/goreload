// Package logger provides structured logging with color support for goreload.
package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
)

// Level represents a logging level.
type Level int

// Log levels.
const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// Logger provides structured logging with optional color and timestamp support.
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	// SetOutput changes the output writer.
	SetOutput(w io.Writer)
	// SetLevel changes the minimum log level.
	SetLevel(level Level)
}

// Config holds logger configuration.
type Config struct {
	Color bool
	Time  bool
	Level string
}

type logger struct {
	mu       sync.Mutex
	out      io.Writer
	level    Level
	useColor bool
	showTime bool

	// Color functions for each level.
	debugColor func(format string, a ...interface{}) string
	infoColor  func(format string, a ...interface{}) string
	warnColor  func(format string, a ...interface{}) string
	errorColor func(format string, a ...interface{}) string
	timeColor  func(format string, a ...interface{}) string
}

// New creates a new Logger with the given configuration.
func New(cfg Config) Logger {
	l := &logger{
		out:      os.Stdout,
		level:    parseLevel(cfg.Level),
		useColor: cfg.Color,
		showTime: cfg.Time,
	}
	l.initColors()
	return l
}

func (l *logger) initColors() {
	if l.useColor {
		l.debugColor = color.New(color.FgHiBlack).SprintfFunc()
		l.infoColor = color.New(color.FgCyan).SprintfFunc()
		l.warnColor = color.New(color.FgYellow).SprintfFunc()
		l.errorColor = color.New(color.FgRed).SprintfFunc()
		l.timeColor = color.New(color.FgHiBlack).SprintfFunc()
	} else {
		noColor := func(format string, a ...interface{}) string {
			return fmt.Sprintf(format, a...)
		}
		l.debugColor = noColor
		l.infoColor = noColor
		l.warnColor = noColor
		l.errorColor = noColor
		l.timeColor = noColor
	}
}

func parseLevel(s string) Level {
	switch s {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	default:
		return LevelInfo
	}
}

// ParseLevel converts a string to a Level.
func ParseLevel(s string) Level {
	return parseLevel(s)
}

func (l *logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = w
}

func (l *logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

func (l *logger) Debug(msg string, args ...any) {
	if l.level > LevelDebug {
		return
	}
	l.log(l.debugColor("[DEBUG]"), msg, args...)
}

func (l *logger) Info(msg string, args ...any) {
	if l.level > LevelInfo {
		return
	}
	l.log(l.infoColor("[INFO]"), msg, args...)
}

func (l *logger) Warn(msg string, args ...any) {
	if l.level > LevelWarn {
		return
	}
	l.log(l.warnColor("[WARN]"), msg, args...)
}

func (l *logger) Error(msg string, args ...any) {
	if l.level > LevelError {
		return
	}
	l.log(l.errorColor("[ERROR]"), msg, args...)
}

func (l *logger) log(prefix, msg string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	var line string
	if l.showTime {
		ts := l.timeColor(time.Now().Format("15:04:05"))
		line = fmt.Sprintf("%s %s %s", ts, prefix, fmt.Sprintf(msg, args...))
	} else {
		line = fmt.Sprintf("%s %s", prefix, fmt.Sprintf(msg, args...))
	}
	_, _ = fmt.Fprintln(l.out, line)
}

// Success prints a success message with a checkmark.
func Success(l Logger, msg string, args ...any) {
	formatted := fmt.Sprintf(msg, args...)
	l.Info("✓ %s", formatted)
}

// Failure prints a failure message with an X mark.
func Failure(l Logger, msg string, args ...any) {
	formatted := fmt.Sprintf(msg, args...)
	l.Error("✗ %s", formatted)
}

// Banner prints the goreload ASCII art banner.
func Banner(w io.Writer, version string) {
	banner := `
   __ _  ___  _ __ ___| | ___   __ _  __| |
  / _` + "`" + ` |/ _ \| '__/ _ \ |/ _ \ / _` + "`" + ` |/ _` + "`" + ` |
 | (_| | (_) | | |  __/ | (_) | (_| | (_| |
  \__, |\___/|_|  \___|_|\___/ \__,_|\__,_|
  |___/                            %s
`
	_, _ = fmt.Fprintf(w, banner, version)
	_, _ = fmt.Fprintln(w)
}
