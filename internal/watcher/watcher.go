// Package watcher provides file system watching with filtering support for goreload.
package watcher

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Op describes a set of file operations.
type Op uint32

// Operations.
const (
	OpCreate Op = 1 << iota
	OpWrite
	OpRemove
	OpRename
	OpChmod
)

func (o Op) String() string {
	switch o {
	case OpCreate:
		return "CREATE"
	case OpWrite:
		return "WRITE"
	case OpRemove:
		return "REMOVE"
	case OpRename:
		return "RENAME"
	case OpChmod:
		return "CHMOD"
	default:
		return "UNKNOWN"
	}
}

// Event represents a file system event.
type Event struct {
	Path string
	Op   Op
	Time time.Time
}

// Watcher watches directories for file changes.
type Watcher interface {
	// Start begins watching and processing events.
	Start(ctx context.Context) error
	// Events returns the channel of filtered events.
	Events() <-chan Event
	// Errors returns the channel of errors.
	Errors() <-chan error
	// Close stops watching and releases resources.
	Close() error
}

// Config holds watcher configuration.
type Config struct {
	Dirs         []string
	Filter       Filter
	Debounce     time.Duration
	Root         string
	ExcludeDirs  []string
}

type watcher struct {
	cfg     Config
	fw      *fsnotify.Watcher
	events  chan Event
	errors  chan error
	done    chan struct{}
	mu      sync.Mutex
	started bool
}

// New creates a new Watcher with the given configuration.
func New(cfg Config) (Watcher, error) {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("create fsnotify watcher: %w", err)
	}

	return &watcher{
		cfg:    cfg,
		fw:     fw,
		events: make(chan Event, 100),
		errors: make(chan error, 10),
		done:   make(chan struct{}),
	}, nil
}

func (w *watcher) Start(ctx context.Context) error {
	w.mu.Lock()
	if w.started {
		w.mu.Unlock()
		return fmt.Errorf("watcher already started")
	}
	w.started = true
	w.mu.Unlock()

	// Add directories to watch.
	for _, dir := range w.cfg.Dirs {
		if !filepath.IsAbs(dir) && w.cfg.Root != "" {
			dir = filepath.Join(w.cfg.Root, dir)
		}

		if err := w.addRecursive(dir); err != nil {
			return fmt.Errorf("add watch directory %s: %w", dir, err)
		}
	}

	go w.loop(ctx)

	return nil
}

func (w *watcher) addRecursive(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		// Skip excluded directories.
		if w.isExcludedDir(path) {
			return filepath.SkipDir
		}

		if err := w.fw.Add(path); err != nil {
			return fmt.Errorf("add %s to watcher: %w", path, err)
		}

		return nil
	})
}

func (w *watcher) isExcludedDir(path string) bool {
	base := filepath.Base(path)
	for _, excluded := range w.cfg.ExcludeDirs {
		if base == excluded {
			return true
		}
	}
	return false
}

func (w *watcher) loop(ctx context.Context) {
	var (
		debounceTimer *time.Timer
		pendingEvents = make(map[string]Event)
		mu            sync.Mutex
	)

	debounce := w.cfg.Debounce
	if debounce <= 0 {
		debounce = 100 * time.Millisecond
	}

	flushEvents := func() {
		mu.Lock()
		defer mu.Unlock()

		for _, evt := range pendingEvents {
			select {
			case w.events <- evt:
			default:
				// Channel full, drop oldest events.
			}
		}
		pendingEvents = make(map[string]Event)
	}

	for {
		select {
		case <-ctx.Done():
			if debounceTimer != nil {
				debounceTimer.Stop()
			}
			flushEvents()
			return

		case <-w.done:
			if debounceTimer != nil {
				debounceTimer.Stop()
			}
			return

		case event, ok := <-w.fw.Events:
			if !ok {
				return
			}

			// Skip if filter doesn't match.
			if w.cfg.Filter != nil && !w.cfg.Filter.Match(event.Name) {
				continue
			}

			// Handle directory creation - add to watch.
			if event.Has(fsnotify.Create) {
				if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
					if !w.isExcludedDir(event.Name) {
						_ = w.fw.Add(event.Name)
					}
				}
			}

			// Convert fsnotify operation to our Op type.
			op := convertOp(event.Op)

			mu.Lock()
			pendingEvents[event.Name] = Event{
				Path: event.Name,
				Op:   op,
				Time: time.Now(),
			}

			if debounceTimer == nil {
				debounceTimer = time.AfterFunc(debounce, flushEvents)
			} else {
				debounceTimer.Reset(debounce)
			}
			mu.Unlock()

		case err, ok := <-w.fw.Errors:
			if !ok {
				return
			}
			select {
			case w.errors <- err:
			default:
			}
		}
	}
}

func convertOp(op fsnotify.Op) Op {
	switch {
	case op.Has(fsnotify.Create):
		return OpCreate
	case op.Has(fsnotify.Write):
		return OpWrite
	case op.Has(fsnotify.Remove):
		return OpRemove
	case op.Has(fsnotify.Rename):
		return OpRename
	case op.Has(fsnotify.Chmod):
		return OpChmod
	default:
		return OpWrite
	}
}

func (w *watcher) Events() <-chan Event {
	return w.events
}

func (w *watcher) Errors() <-chan error {
	return w.errors
}

func (w *watcher) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	select {
	case <-w.done:
		// Already closed.
	default:
		close(w.done)
	}

	return w.fw.Close()
}
