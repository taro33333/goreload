// Package runner provides process management for goreload.
package runner

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

// Runner manages the lifecycle of the target application process.
type Runner interface {
	// Start launches the application process.
	Start(ctx context.Context) error
	// Stop terminates the application process gracefully.
	Stop(ctx context.Context) error
	// Restart stops and then starts the application process.
	Restart(ctx context.Context) error
	// Running returns true if the application is currently running.
	Running() bool
}

// Config holds runner configuration.
type Config struct {
	Bin       string
	Args      []string
	Root      string
	KillDelay time.Duration
	Stdout    io.Writer
	Stderr    io.Writer
}

type runner struct {
	cfg Config

	mu      sync.Mutex
	cmd     *exec.Cmd
	running bool
	done    chan struct{}
}

// New creates a new Runner with the given configuration.
func New(cfg Config) Runner {
	if cfg.Stdout == nil {
		cfg.Stdout = os.Stdout
	}
	if cfg.Stderr == nil {
		cfg.Stderr = os.Stderr
	}
	return &runner{cfg: cfg}
}

func (r *runner) Start(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.running {
		return fmt.Errorf("process already running")
	}

	binPath := r.cfg.Bin
	if !filepath.IsAbs(binPath) {
		binPath = filepath.Join(r.cfg.Root, binPath)
	}

	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		return fmt.Errorf("binary not found: %s", binPath)
	}

	r.cmd = exec.CommandContext(ctx, binPath, r.cfg.Args...)
	r.cmd.Dir = r.cfg.Root
	r.cmd.Stdout = r.cfg.Stdout
	r.cmd.Stderr = r.cfg.Stderr

	prepareCommand(r.cmd)

	if err := r.cmd.Start(); err != nil {
		return fmt.Errorf("start process: %w", err)
	}

	r.running = true
	r.done = make(chan struct{})

	go r.wait()

	return nil
}

func (r *runner) wait() {
	if r.cmd != nil && r.cmd.Process != nil {
		_ = r.cmd.Wait()
	}

	r.mu.Lock()
	r.running = false
	if r.done != nil {
		close(r.done)
	}
	r.mu.Unlock()
}

func (r *runner) Stop(ctx context.Context) error {
	r.mu.Lock()
	if !r.running || r.cmd == nil || r.cmd.Process == nil {
		r.mu.Unlock()
		return nil
	}

	done := r.done
	proc := r.cmd.Process
	r.mu.Unlock()

	// First try graceful shutdown.
	_ = interruptProcess(proc)

	// Wait for graceful shutdown or timeout.
	killDelay := r.cfg.KillDelay
	if killDelay <= 0 {
		killDelay = 500 * time.Millisecond
	}

	select {
	case <-done:
		return nil
	case <-time.After(killDelay):
		// Force kill.
		_ = killProcess(proc)
	case <-ctx.Done():
		// Force kill on context cancellation.
		_ = killProcess(proc)
		return ctx.Err()
	}

	// Wait for process to actually exit.
	select {
	case <-done:
		return nil
	case <-time.After(time.Second):
		return fmt.Errorf("process did not exit after SIGKILL")
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r *runner) Restart(ctx context.Context) error {
	if err := r.Stop(ctx); err != nil {
		return fmt.Errorf("stop for restart: %w", err)
	}
	if err := r.Start(ctx); err != nil {
		return fmt.Errorf("start for restart: %w", err)
	}
	return nil
}

func (r *runner) Running() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.running
}
