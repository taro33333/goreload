// Package config provides configuration structures and validation for goreload.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Default configuration values.
const (
	DefaultTmpDir     = "tmp"
	DefaultBuildCmd   = "go build -o ./tmp/main ."
	DefaultBin        = "./tmp/main"
	DefaultDelay      = 200 * time.Millisecond
	DefaultKillDelay  = 500 * time.Millisecond
	DefaultLogLevel   = "info"
	DefaultConfigFile = "goreload.yaml"
)

// Sentinel errors for configuration validation.
var (
	ErrEmptyBuildCmd    = errors.New("build command cannot be empty")
	ErrEmptyBin         = errors.New("binary path cannot be empty")
	ErrInvalidDelay     = errors.New("delay must be positive")
	ErrInvalidKillDelay = errors.New("kill_delay must be positive")
	ErrInvalidLogLevel  = errors.New("log level must be one of: debug, info, warn, error")
	ErrNoExtensions     = errors.New("at least one file extension must be specified")
	ErrNoDirs           = errors.New("at least one watch directory must be specified")
)

// Config represents the complete goreload configuration.
type Config struct {
	Root   string      `yaml:"root"`
	TmpDir string      `yaml:"tmp_dir"`
	Build  BuildConfig `yaml:"build"`
	Watch  WatchConfig `yaml:"watch"`
	Log    LogConfig   `yaml:"log"`
}

// BuildConfig holds build-related settings.
type BuildConfig struct {
	Cmd       string        `yaml:"cmd"`
	Bin       string        `yaml:"bin"`
	Args      []string      `yaml:"args"`
	Delay     time.Duration `yaml:"delay"`
	KillDelay time.Duration `yaml:"kill_delay"`
}

// WatchConfig holds file watching settings.
type WatchConfig struct {
	Extensions   []string `yaml:"extensions"`
	Dirs         []string `yaml:"dirs"`
	ExcludeDirs  []string `yaml:"exclude_dirs"`
	ExcludeFiles []string `yaml:"exclude_files"`
}

// LogConfig holds logging settings.
type LogConfig struct {
	Color bool   `yaml:"color"`
	Time  bool   `yaml:"time"`
	Level string `yaml:"level"`
}

// Validate checks all configuration values and returns an error if any are invalid.
func (c *Config) Validate() error {
	if err := c.Build.validate(); err != nil {
		return fmt.Errorf("build config: %w", err)
	}
	if err := c.Watch.validate(); err != nil {
		return fmt.Errorf("watch config: %w", err)
	}
	if err := c.Log.validate(); err != nil {
		return fmt.Errorf("log config: %w", err)
	}
	return nil
}

func (b *BuildConfig) validate() error {
	if b.Cmd == "" {
		return ErrEmptyBuildCmd
	}
	if b.Bin == "" {
		return ErrEmptyBin
	}
	if b.Delay < 0 {
		return ErrInvalidDelay
	}
	if b.KillDelay < 0 {
		return ErrInvalidKillDelay
	}
	return nil
}

func (w *WatchConfig) validate() error {
	if len(w.Extensions) == 0 {
		return ErrNoExtensions
	}
	if len(w.Dirs) == 0 {
		return ErrNoDirs
	}
	return nil
}

func (l *LogConfig) validate() error {
	switch l.Level {
	case "debug", "info", "warn", "error":
		return nil
	default:
		return ErrInvalidLogLevel
	}
}

// AbsRoot returns the absolute path of the root directory.
func (c *Config) AbsRoot() (string, error) {
	if c.Root == "" || c.Root == "." {
		return os.Getwd()
	}
	return filepath.Abs(c.Root)
}

// AbsTmpDir returns the absolute path of the temporary directory.
func (c *Config) AbsTmpDir() (string, error) {
	root, err := c.AbsRoot()
	if err != nil {
		return "", err
	}
	if filepath.IsAbs(c.TmpDir) {
		return c.TmpDir, nil
	}
	return filepath.Join(root, c.TmpDir), nil
}
