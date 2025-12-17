// Package builder provides build command execution for goreload.
package builder

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Result contains the outcome of a build operation.
type Result struct {
	Success  bool
	Output   string
	Duration time.Duration
	Error    error
}

// Builder executes build commands and manages build artifacts.
type Builder interface {
	// Build executes the build command and returns the result.
	Build(ctx context.Context) Result
	// Clean removes build artifacts.
	Clean() error
}

// Config holds builder configuration.
type Config struct {
	Cmd    string
	Bin    string
	TmpDir string
	Root   string
}

type builder struct {
	cfg Config
}

// New creates a new Builder with the given configuration.
func New(cfg Config) Builder {
	return &builder{cfg: cfg}
}

func (b *builder) Build(ctx context.Context) Result {
	start := time.Now()

	if err := b.ensureTmpDir(); err != nil {
		return Result{
			Success:  false,
			Output:   "",
			Duration: time.Since(start),
			Error:    fmt.Errorf("create tmp dir: %w", err),
		}
	}

	args := parseCommand(b.cfg.Cmd)
	if len(args) == 0 {
		return Result{
			Success:  false,
			Output:   "",
			Duration: time.Since(start),
			Error:    fmt.Errorf("empty build command"),
		}
	}

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Dir = b.cfg.Root

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	duration := time.Since(start)

	output := stdout.String()
	if stderr.Len() > 0 {
		if output != "" {
			output += "\n"
		}
		output += stderr.String()
	}

	if err != nil {
		if ctx.Err() == context.Canceled {
			return Result{
				Success:  false,
				Output:   output,
				Duration: duration,
				Error:    ctx.Err(),
			}
		}
		return Result{
			Success:  false,
			Output:   output,
			Duration: duration,
			Error:    fmt.Errorf("build failed: %w", err),
		}
	}

	return Result{
		Success:  true,
		Output:   output,
		Duration: duration,
		Error:    nil,
	}
}

func (b *builder) Clean() error {
	binPath := b.cfg.Bin
	if !filepath.IsAbs(binPath) {
		binPath = filepath.Join(b.cfg.Root, binPath)
	}

	if err := os.Remove(binPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove binary: %w", err)
	}
	return nil
}

func (b *builder) ensureTmpDir() error {
	tmpDir := b.cfg.TmpDir
	if !filepath.IsAbs(tmpDir) {
		tmpDir = filepath.Join(b.cfg.Root, tmpDir)
	}
	return os.MkdirAll(tmpDir, 0755)
}

// parseCommand splits a command string into arguments.
// It handles simple quoting but is not a full shell parser.
func parseCommand(cmd string) []string {
	var args []string
	var current strings.Builder
	inQuote := false
	quoteChar := rune(0)

	for _, r := range cmd {
		switch {
		case r == '"' || r == '\'':
			if inQuote && r == quoteChar {
				inQuote = false
				quoteChar = 0
			} else if !inQuote {
				inQuote = true
				quoteChar = r
			} else {
				current.WriteRune(r)
			}
		case r == ' ' && !inQuote:
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args
}
