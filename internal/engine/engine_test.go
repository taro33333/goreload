package engine

import (
	"context"
	"testing"
	"time"

	"github.com/user/goreload/internal/config"
	"github.com/user/goreload/internal/logger"
)

func TestNew(t *testing.T) {
	cfg := &config.Config{
		Root:   ".",
		TmpDir: "tmp",
		Build: config.BuildConfig{
			Cmd:       "go build -o ./tmp/main .",
			Bin:       "./tmp/main",
			Args:      []string{},
			Delay:     200 * time.Millisecond,
			KillDelay: 500 * time.Millisecond,
		},
		Watch: config.WatchConfig{
			Extensions:   []string{".go"},
			Dirs:         []string{"."},
			ExcludeDirs:  []string{"tmp"},
			ExcludeFiles: []string{},
		},
		Log: config.LogConfig{
			Color: false,
			Time:  false,
			Level: "info",
		},
	}

	log := logger.New(logger.Config{
		Color: false,
		Time:  false,
		Level: "info",
	})

	eng, err := New(cfg, log)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if eng == nil {
		t.Error("New() returned nil")
	}
}

func TestEngine_RunWithContextCancel(t *testing.T) {
	cfg := &config.Config{
		Root:   t.TempDir(),
		TmpDir: "tmp",
		Build: config.BuildConfig{
			Cmd:       "echo test",
			Bin:       "/bin/echo",
			Args:      []string{},
			Delay:     50 * time.Millisecond,
			KillDelay: 100 * time.Millisecond,
		},
		Watch: config.WatchConfig{
			Extensions:   []string{".go"},
			Dirs:         []string{"."},
			ExcludeDirs:  []string{"tmp"},
			ExcludeFiles: []string{},
		},
		Log: config.LogConfig{
			Color: false,
			Time:  false,
			Level: "error",
		},
	}

	log := logger.New(logger.Config{
		Color: false,
		Time:  false,
		Level: "error",
	})

	eng, err := New(cfg, log)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- eng.Run(ctx)
	}()

	select {
	case err := <-done:
		if err != nil && err != context.DeadlineExceeded && err != context.Canceled {
			t.Errorf("Run() error = %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Error("Run() did not return after context cancellation")
	}
}

func TestEngine_DoubleRun(t *testing.T) {
	cfg := &config.Config{
		Root:   t.TempDir(),
		TmpDir: "tmp",
		Build: config.BuildConfig{
			Cmd:       "echo test",
			Bin:       "/bin/echo",
			Args:      []string{},
			Delay:     50 * time.Millisecond,
			KillDelay: 100 * time.Millisecond,
		},
		Watch: config.WatchConfig{
			Extensions:   []string{".go"},
			Dirs:         []string{"."},
			ExcludeDirs:  []string{"tmp"},
			ExcludeFiles: []string{},
		},
		Log: config.LogConfig{
			Color: false,
			Time:  false,
			Level: "error",
		},
	}

	log := logger.New(logger.Config{
		Color: false,
		Time:  false,
		Level: "error",
	})

	eng, err := New(cfg, log)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start first run in background
	go func() {
		_ = eng.Run(ctx)
	}()

	// Wait for engine to start
	time.Sleep(100 * time.Millisecond)

	// Try to run again - should fail
	err = eng.Run(ctx)
	if err == nil {
		t.Error("Run() should fail when already running")
	}

	cancel()
	time.Sleep(200 * time.Millisecond)
}

func TestEngine_Stop(t *testing.T) {
	cfg := &config.Config{
		Root:   t.TempDir(),
		TmpDir: "tmp",
		Build: config.BuildConfig{
			Cmd:       "echo test",
			Bin:       "/bin/echo",
			Args:      []string{},
			Delay:     50 * time.Millisecond,
			KillDelay: 100 * time.Millisecond,
		},
		Watch: config.WatchConfig{
			Extensions:   []string{".go"},
			Dirs:         []string{"."},
			ExcludeDirs:  []string{"tmp"},
			ExcludeFiles: []string{},
		},
		Log: config.LogConfig{
			Color: false,
			Time:  false,
			Level: "error",
		},
	}

	log := logger.New(logger.Config{
		Color: false,
		Time:  false,
		Level: "error",
	})

	eng, err := New(cfg, log)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Stop when not running should not error
	ctx := context.Background()
	if err := eng.Stop(ctx); err != nil {
		t.Errorf("Stop() error = %v when not running", err)
	}
}
