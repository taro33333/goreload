// Package engine provides the main orchestrator for goreload.
package engine

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/taro33333/goreload/internal/builder"
	"github.com/taro33333/goreload/internal/config"
	"github.com/taro33333/goreload/internal/logger"
	"github.com/taro33333/goreload/internal/runner"
	"github.com/taro33333/goreload/internal/watcher"
)

// Engine orchestrates the watch-build-run cycle.
type Engine struct {
	cfg     *config.Config
	log     logger.Logger
	builder builder.Builder
	runner  runner.Runner
	watcher watcher.Watcher

	mu      sync.Mutex
	running bool
}

// New creates a new Engine with the given configuration.
func New(cfg *config.Config, log logger.Logger) (*Engine, error) {
	root, err := cfg.AbsRoot()
	if err != nil {
		return nil, fmt.Errorf("resolve root: %w", err)
	}

	tmpDir, err := cfg.AbsTmpDir()
	if err != nil {
		return nil, fmt.Errorf("resolve tmp dir: %w", err)
	}

	b := builder.New(builder.Config{
		Cmd:    cfg.Build.Cmd,
		Bin:    cfg.Build.Bin,
		TmpDir: tmpDir,
		Root:   root,
	})

	r := runner.New(runner.Config{
		Bin:       cfg.Build.Bin,
		Args:      cfg.Build.Args,
		Root:      root,
		KillDelay: cfg.Build.KillDelay,
	})

	f := watcher.NewFilter(watcher.FilterConfig{
		Extensions:   cfg.Watch.Extensions,
		ExcludeDirs:  cfg.Watch.ExcludeDirs,
		ExcludeFiles: cfg.Watch.ExcludeFiles,
		Root:         root,
	})

	w, err := watcher.New(watcher.Config{
		Dirs:        cfg.Watch.Dirs,
		Filter:      f,
		Debounce:    cfg.Build.Delay,
		Root:        root,
		ExcludeDirs: cfg.Watch.ExcludeDirs,
	})
	if err != nil {
		return nil, fmt.Errorf("create watcher: %w", err)
	}

	return &Engine{
		cfg:     cfg,
		log:     log,
		builder: b,
		runner:  r,
		watcher: w,
	}, nil
}

// Run starts the main engine loop.
func (e *Engine) Run(ctx context.Context) error {
	e.mu.Lock()
	if e.running {
		e.mu.Unlock()
		return fmt.Errorf("engine already running")
	}
	e.running = true
	e.mu.Unlock()

	defer func() {
		e.mu.Lock()
		e.running = false
		e.mu.Unlock()
	}()

	root, _ := e.cfg.AbsRoot()

	// Log watched directories.
	for _, dir := range e.cfg.Watch.Dirs {
		if !filepath.IsAbs(dir) {
			dir = filepath.Join(root, dir)
		}
		e.log.Info("watching: %s", dir)
	}

	// Log excluded directories.
	if len(e.cfg.Watch.ExcludeDirs) > 0 {
		e.log.Info("excluding: %v", e.cfg.Watch.ExcludeDirs)
	}

	// Start watcher.
	if err := e.watcher.Start(ctx); err != nil {
		return fmt.Errorf("start watcher: %w", err)
	}
	defer func() { _ = e.watcher.Close() }()

	// Initial build and run.
	if err := e.buildAndRun(ctx); err != nil {
		e.log.Error("initial build failed: %v", err)
		// Continue watching for changes even if initial build fails.
	}

	// Main loop.
	for {
		select {
		case <-ctx.Done():
			e.log.Info("shutting down...")
			stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			_ = e.runner.Stop(stopCtx)
			cancel()
			return ctx.Err()

		case evt, ok := <-e.watcher.Events():
			if !ok {
				return nil
			}

			relPath, err := filepath.Rel(root, evt.Path)
			if err != nil {
				relPath = evt.Path
			}
			e.log.Info("%s changed", relPath)

			if err := e.buildAndRun(ctx); err != nil {
				e.log.Error("rebuild failed: %v", err)
			}

		case err, ok := <-e.watcher.Errors():
			if !ok {
				return nil
			}
			e.log.Error("watcher error: %v", err)
		}
	}
}

func (e *Engine) buildAndRun(ctx context.Context) error {
	// Stop current process.
	if e.runner.Running() {
		stopCtx, cancel := context.WithTimeout(ctx, e.cfg.Build.KillDelay*2)
		if err := e.runner.Stop(stopCtx); err != nil {
			cancel()
			e.log.Warn("failed to stop process: %v", err)
		}
		cancel()
	}

	// Build.
	e.log.Info("building...")
	result := e.builder.Build(ctx)

	if !result.Success {
		if result.Output != "" {
			e.log.Error("build output:\n%s", result.Output)
		}
		logger.Failure(e.log, "build failed (%.2fs)", result.Duration.Seconds())
		return result.Error
	}

	logger.Success(e.log, "build completed (%.2fs)", result.Duration.Seconds())

	// Run.
	if err := e.runner.Start(ctx); err != nil {
		logger.Failure(e.log, "failed to start: %v", err)
		return err
	}

	logger.Success(e.log, "running %s", e.cfg.Build.Bin)
	return nil
}

// Stop gracefully stops the engine.
func (e *Engine) Stop(ctx context.Context) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.running {
		return nil
	}

	if err := e.runner.Stop(ctx); err != nil {
		return fmt.Errorf("stop runner: %w", err)
	}

	if err := e.watcher.Close(); err != nil {
		return fmt.Errorf("close watcher: %w", err)
	}

	return nil
}
