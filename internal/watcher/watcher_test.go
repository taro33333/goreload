package watcher

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	cfg := Config{
		Dirs:        []string{"."},
		Filter:      nil,
		Debounce:    100 * time.Millisecond,
		Root:        ".",
		ExcludeDirs: []string{"tmp"},
	}

	w, err := New(cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer w.Close()

	if w == nil {
		t.Error("New() returned nil")
	}
}

func TestWatcher_StartAndClose(t *testing.T) {
	tmpDir := t.TempDir()

	w, err := New(Config{
		Dirs:        []string{"."},
		Debounce:    50 * time.Millisecond,
		Root:        tmpDir,
		ExcludeDirs: []string{},
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := w.Start(ctx); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	// Starting again should fail
	if err := w.Start(ctx); err == nil {
		t.Error("Start() should fail when already started")
	}

	if err := w.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

func TestWatcher_Events(t *testing.T) {
	tmpDir := t.TempDir()

	filter := NewFilter(FilterConfig{
		Extensions:   []string{".go"},
		ExcludeDirs:  []string{},
		ExcludeFiles: []string{},
		Root:         tmpDir,
	})

	w, err := New(Config{
		Dirs:        []string{"."},
		Filter:      filter,
		Debounce:    50 * time.Millisecond,
		Root:        tmpDir,
		ExcludeDirs: []string{},
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer w.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := w.Start(ctx); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	// Create a file
	goFile := filepath.Join(tmpDir, "test.go")
	if err := os.WriteFile(goFile, []byte("package main"), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	// Wait for event
	select {
	case evt := <-w.Events():
		if evt.Path != goFile {
			t.Errorf("Event path = %v, want %v", evt.Path, goFile)
		}
	case err := <-w.Errors():
		t.Errorf("Unexpected error: %v", err)
	case <-time.After(500 * time.Millisecond):
		t.Error("Timeout waiting for event")
	}
}

func TestWatcher_FilteredEvents(t *testing.T) {
	tmpDir := t.TempDir()

	filter := NewFilter(FilterConfig{
		Extensions:   []string{".go"},
		ExcludeDirs:  []string{},
		ExcludeFiles: []string{"*_test.go"},
		Root:         tmpDir,
	})

	w, err := New(Config{
		Dirs:        []string{"."},
		Filter:      filter,
		Debounce:    50 * time.Millisecond,
		Root:        tmpDir,
		ExcludeDirs: []string{},
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer w.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := w.Start(ctx); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	// Create a test file (should be filtered)
	testFile := filepath.Join(tmpDir, "main_test.go")
	if err := os.WriteFile(testFile, []byte("package main"), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	// Should not receive event
	select {
	case evt := <-w.Events():
		t.Errorf("Should not receive event for test file, got: %v", evt)
	case <-time.After(200 * time.Millisecond):
		// Expected - no event
	}

	// Create a regular file (should trigger event)
	goFile := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(goFile, []byte("package main"), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	select {
	case evt := <-w.Events():
		if evt.Path != goFile {
			t.Errorf("Event path = %v, want %v", evt.Path, goFile)
		}
	case <-time.After(500 * time.Millisecond):
		t.Error("Timeout waiting for event")
	}
}

func TestWatcher_ExcludeDirs(t *testing.T) {
	tmpDir := t.TempDir()

	// Create excluded directory
	excludedDir := filepath.Join(tmpDir, "vendor")
	if err := os.MkdirAll(excludedDir, 0755); err != nil {
		t.Fatalf("create vendor dir: %v", err)
	}

	w, err := New(Config{
		Dirs:        []string{"."},
		Debounce:    50 * time.Millisecond,
		Root:        tmpDir,
		ExcludeDirs: []string{"vendor"},
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer w.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := w.Start(ctx); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	// Create file in excluded directory
	vendorFile := filepath.Join(excludedDir, "lib.go")
	if err := os.WriteFile(vendorFile, []byte("package lib"), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	// Should not receive event
	select {
	case evt := <-w.Events():
		t.Errorf("Should not receive event for vendor file, got: %v", evt)
	case <-time.After(200 * time.Millisecond):
		// Expected - no event
	}
}

func TestWatcher_ContextCancellation(t *testing.T) {
	tmpDir := t.TempDir()

	w, err := New(Config{
		Dirs:        []string{"."},
		Debounce:    50 * time.Millisecond,
		Root:        tmpDir,
		ExcludeDirs: []string{},
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer w.Close()

	ctx, cancel := context.WithCancel(context.Background())

	if err := w.Start(ctx); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	// Cancel context
	cancel()

	// Give time for cancellation to propagate
	time.Sleep(100 * time.Millisecond)

	// Should be able to close without hanging
	done := make(chan struct{})
	go func() {
		w.Close()
		close(done)
	}()

	select {
	case <-done:
		// Success
	case <-time.After(time.Second):
		t.Error("Close() hung after context cancellation")
	}
}

func TestOp_String(t *testing.T) {
	tests := []struct {
		op   Op
		want string
	}{
		{OpCreate, "CREATE"},
		{OpWrite, "WRITE"},
		{OpRemove, "REMOVE"},
		{OpRename, "RENAME"},
		{OpChmod, "CHMOD"},
		{Op(999), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.op.String(); got != tt.want {
				t.Errorf("Op.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
