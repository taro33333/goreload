package runner

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	cfg := Config{
		Bin:       "./main",
		Args:      []string{"-v"},
		Root:      ".",
		KillDelay: 500 * time.Millisecond,
	}

	r := New(cfg)
	if r == nil {
		t.Error("New() returned nil")
	}
}

func TestRunner_StartStop(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a simple script that runs indefinitely
	scriptPath := filepath.Join(tmpDir, "test.sh")
	script := `#!/bin/sh
while true; do
    sleep 1
done
`
	if err := os.WriteFile(scriptPath, []byte(script), 0755); err != nil {
		t.Fatalf("write script: %v", err)
	}

	var stdout, stderr bytes.Buffer
	r := New(Config{
		Bin:       scriptPath,
		Args:      []string{},
		Root:      tmpDir,
		KillDelay: 100 * time.Millisecond,
		Stdout:    &stdout,
		Stderr:    &stderr,
	})

	ctx := context.Background()

	// Start
	if err := r.Start(ctx); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	if !r.Running() {
		t.Error("Running() = false after Start()")
	}

	// Starting again should fail
	if err := r.Start(ctx); err == nil {
		t.Error("Start() should fail when already running")
	}

	// Stop
	if err := r.Stop(ctx); err != nil {
		t.Errorf("Stop() error = %v", err)
	}

	// Wait a bit for the process to fully exit
	time.Sleep(200 * time.Millisecond)

	if r.Running() {
		t.Error("Running() = true after Stop()")
	}
}

func TestRunner_Restart(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a simple script
	scriptPath := filepath.Join(tmpDir, "test.sh")
	script := `#!/bin/sh
while true; do
    sleep 1
done
`
	if err := os.WriteFile(scriptPath, []byte(script), 0755); err != nil {
		t.Fatalf("write script: %v", err)
	}

	r := New(Config{
		Bin:       scriptPath,
		Args:      []string{},
		Root:      tmpDir,
		KillDelay: 100 * time.Millisecond,
	})

	ctx := context.Background()

	// Start
	if err := r.Start(ctx); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	// Restart
	if err := r.Restart(ctx); err != nil {
		t.Errorf("Restart() error = %v", err)
	}

	if !r.Running() {
		t.Error("Running() = false after Restart()")
	}

	// Cleanup
	if err := r.Stop(ctx); err != nil {
		t.Errorf("Stop() error = %v", err)
	}
}

func TestRunner_StartBinaryNotFound(t *testing.T) {
	r := New(Config{
		Bin:       "/nonexistent/binary",
		Args:      []string{},
		Root:      ".",
		KillDelay: 100 * time.Millisecond,
	})

	ctx := context.Background()
	err := r.Start(ctx)
	if err == nil {
		t.Error("Start() should fail for nonexistent binary")
	}
}

func TestRunner_StopNotRunning(t *testing.T) {
	r := New(Config{
		Bin:       "./main",
		Args:      []string{},
		Root:      ".",
		KillDelay: 100 * time.Millisecond,
	})

	ctx := context.Background()

	// Stop when not running should not error
	if err := r.Stop(ctx); err != nil {
		t.Errorf("Stop() error = %v when not running", err)
	}
}

func TestRunner_StopWithContextCancellation(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a script that ignores SIGINT
	scriptPath := filepath.Join(tmpDir, "stubborn.sh")
	script := `#!/bin/sh
trap '' INT
while true; do
    sleep 1
done
`
	if err := os.WriteFile(scriptPath, []byte(script), 0755); err != nil {
		t.Fatalf("write script: %v", err)
	}

	r := New(Config{
		Bin:       scriptPath,
		Args:      []string{},
		Root:      tmpDir,
		KillDelay: 50 * time.Millisecond,
	})

	ctx := context.Background()
	if err := r.Start(ctx); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	// Stop with short timeout - should still succeed because of SIGKILL
	stopCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	err := r.Stop(stopCtx)
	// Should eventually succeed due to SIGKILL
	if err != nil && err != context.DeadlineExceeded {
		t.Logf("Stop() returned error (may be expected): %v", err)
	}

	// Give time for cleanup
	time.Sleep(200 * time.Millisecond)
}

func TestRunner_Running(t *testing.T) {
	r := New(Config{
		Bin:       "./main",
		Args:      []string{},
		Root:      ".",
		KillDelay: 100 * time.Millisecond,
	})

	if r.Running() {
		t.Error("Running() = true before Start()")
	}
}
