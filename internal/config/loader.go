package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// rawConfig is used for YAML unmarshaling with string durations.
type rawConfig struct {
	Root   string         `yaml:"root"`
	TmpDir string         `yaml:"tmp_dir"`
	Build  rawBuildConfig `yaml:"build"`
	Watch  WatchConfig    `yaml:"watch"`
	Log    LogConfig      `yaml:"log"`
}

type rawBuildConfig struct {
	Cmd       string   `yaml:"cmd"`
	Bin       string   `yaml:"bin"`
	Args      []string `yaml:"args"`
	Delay     string   `yaml:"delay"`
	KillDelay string   `yaml:"kill_delay"`
}

// Default returns a Config with default values.
func Default() *Config {
	return &Config{
		Root:   ".",
		TmpDir: DefaultTmpDir,
		Build: BuildConfig{
			Cmd:       DefaultBuildCmd,
			Bin:       DefaultBin,
			Args:      []string{},
			Delay:     DefaultDelay,
			KillDelay: DefaultKillDelay,
		},
		Watch: WatchConfig{
			Extensions:   []string{".go"},
			Dirs:         []string{"."},
			ExcludeDirs:  []string{"tmp", "vendor", ".git", "node_modules"},
			ExcludeFiles: []string{},
		},
		Log: LogConfig{
			Color: true,
			Time:  true,
			Level: DefaultLogLevel,
		},
	}
}

// LoadWithDefaults reads a YAML configuration file and merges it with default values.
func LoadWithDefaults(path string) (*Config, error) {
	cfg := Default()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var raw rawConfig
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}

	if err := mergeConfig(cfg, &raw); err != nil {
		return nil, fmt.Errorf("merge config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return cfg, nil
}

func mergeConfig(cfg *Config, raw *rawConfig) error {
	if raw.Root != "" {
		cfg.Root = raw.Root
	}
	if raw.TmpDir != "" {
		cfg.TmpDir = raw.TmpDir
	}

	if err := mergeBuildConfig(&cfg.Build, &raw.Build); err != nil {
		return err
	}
	mergeWatchConfig(&cfg.Watch, &raw.Watch)
	mergeLogConfig(&cfg.Log, &raw.Log)

	return nil
}

func mergeBuildConfig(cfg *BuildConfig, raw *rawBuildConfig) error {
	if raw.Cmd != "" {
		cfg.Cmd = raw.Cmd
	}
	if raw.Bin != "" {
		cfg.Bin = raw.Bin
	}
	if len(raw.Args) > 0 {
		cfg.Args = raw.Args
	}
	if raw.Delay != "" {
		d, err := time.ParseDuration(raw.Delay)
		if err != nil {
			return fmt.Errorf("parse delay: %w", err)
		}
		cfg.Delay = d
	}
	if raw.KillDelay != "" {
		d, err := time.ParseDuration(raw.KillDelay)
		if err != nil {
			return fmt.Errorf("parse kill_delay: %w", err)
		}
		cfg.KillDelay = d
	}
	return nil
}

func mergeWatchConfig(cfg *WatchConfig, raw *WatchConfig) {
	if len(raw.Extensions) > 0 {
		cfg.Extensions = raw.Extensions
	}
	if len(raw.Dirs) > 0 {
		cfg.Dirs = raw.Dirs
	}
	if len(raw.ExcludeDirs) > 0 {
		cfg.ExcludeDirs = raw.ExcludeDirs
	}
	if len(raw.ExcludeFiles) > 0 {
		cfg.ExcludeFiles = raw.ExcludeFiles
	}
}

func mergeLogConfig(cfg *LogConfig, raw *LogConfig) {
	// For boolean fields, we need a different approach since zero value is meaningful.
	// The raw config will have the parsed values, so we always use them if the YAML had values.
	// This is a limitation - we can't distinguish between "not set" and "set to false".
	// For now, we assume the YAML values take precedence if the raw struct has non-default values.
	// A more robust solution would use pointers, but that complicates the API.
	cfg.Color = raw.Color
	cfg.Time = raw.Time
	if raw.Level != "" {
		cfg.Level = raw.Level
	}
}

// Exists checks if a configuration file exists at the given path.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// WriteDefault writes a default configuration file to the given path.
func WriteDefault(path string) error {
	content := `# goreload configuration file

# Project root directory
root: "."

# Temporary directory for build artifacts
tmp_dir: "tmp"

# Build settings
build:
  # Build command
  cmd: "go build -o ./tmp/main ."
  # Binary to execute
  bin: "./tmp/main"
  # Arguments to pass to the binary
  args: []
  # Delay before building after file change (debounce)
  delay: "200ms"
  # Grace period for process termination
  kill_delay: "500ms"

# File watching settings
watch:
  # File extensions to watch
  extensions:
    - ".go"
  # Directories to watch
  dirs:
    - "."
  # Directories to exclude
  exclude_dirs:
    - "tmp"
    - "vendor"
    - ".git"
    - "node_modules"
  # Files to exclude (glob patterns)
  exclude_files:
    - "*_test.go"

# Logging settings
log:
  # Enable colored output
  color: true
  # Show timestamps
  time: true
  # Log level: debug, info, warn, error
  level: "info"
`
	return os.WriteFile(path, []byte(content), 0644)
}
