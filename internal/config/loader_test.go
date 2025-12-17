package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg.Root != "." {
		t.Errorf("Default().Root = %v, want .", cfg.Root)
	}
	if cfg.TmpDir != DefaultTmpDir {
		t.Errorf("Default().TmpDir = %v, want %v", cfg.TmpDir, DefaultTmpDir)
	}
	if cfg.Build.Cmd != DefaultBuildCmd {
		t.Errorf("Default().Build.Cmd = %v, want %v", cfg.Build.Cmd, DefaultBuildCmd)
	}
	if cfg.Build.Bin != DefaultBin {
		t.Errorf("Default().Build.Bin = %v, want %v", cfg.Build.Bin, DefaultBin)
	}
	if cfg.Build.Delay != DefaultDelay {
		t.Errorf("Default().Build.Delay = %v, want %v", cfg.Build.Delay, DefaultDelay)
	}
	if cfg.Build.KillDelay != DefaultKillDelay {
		t.Errorf("Default().Build.KillDelay = %v, want %v", cfg.Build.KillDelay, DefaultKillDelay)
	}
	if cfg.Log.Level != DefaultLogLevel {
		t.Errorf("Default().Log.Level = %v, want %v", cfg.Log.Level, DefaultLogLevel)
	}
	if len(cfg.Watch.Extensions) == 0 || cfg.Watch.Extensions[0] != ".go" {
		t.Errorf("Default().Watch.Extensions = %v, want [.go]", cfg.Watch.Extensions)
	}
}

func TestLoadWithDefaults(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("valid config file", func(t *testing.T) {
		configPath := filepath.Join(tmpDir, "valid.yaml")
		content := `
root: "./myapp"
tmp_dir: "build"
build:
  cmd: "go build -o ./build/app ."
  bin: "./build/app"
  delay: "300ms"
  kill_delay: "1s"
watch:
  extensions:
    - ".go"
    - ".html"
  dirs:
    - "."
    - "templates"
  exclude_dirs:
    - "build"
log:
  level: "debug"
`
		if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
			t.Fatalf("write config file: %v", err)
		}

		cfg, err := LoadWithDefaults(configPath)
		if err != nil {
			t.Fatalf("LoadWithDefaults() error = %v", err)
		}

		if cfg.Root != "./myapp" {
			t.Errorf("Root = %v, want ./myapp", cfg.Root)
		}
		if cfg.TmpDir != "build" {
			t.Errorf("TmpDir = %v, want build", cfg.TmpDir)
		}
		if cfg.Build.Cmd != "go build -o ./build/app ." {
			t.Errorf("Build.Cmd = %v", cfg.Build.Cmd)
		}
		if cfg.Build.Delay != 300*time.Millisecond {
			t.Errorf("Build.Delay = %v, want 300ms", cfg.Build.Delay)
		}
		if cfg.Build.KillDelay != 1*time.Second {
			t.Errorf("Build.KillDelay = %v, want 1s", cfg.Build.KillDelay)
		}
		if len(cfg.Watch.Extensions) != 2 {
			t.Errorf("Watch.Extensions = %v, want 2 items", cfg.Watch.Extensions)
		}
		if cfg.Log.Level != "debug" {
			t.Errorf("Log.Level = %v, want debug", cfg.Log.Level)
		}
	})

	t.Run("file not found", func(t *testing.T) {
		_, err := LoadWithDefaults(filepath.Join(tmpDir, "nonexistent.yaml"))
		if err == nil {
			t.Error("LoadWithDefaults() error = nil, want error")
		}
	})

	t.Run("invalid yaml", func(t *testing.T) {
		configPath := filepath.Join(tmpDir, "invalid.yaml")
		content := `
root: [invalid
`
		if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
			t.Fatalf("write config file: %v", err)
		}

		_, err := LoadWithDefaults(configPath)
		if err == nil {
			t.Error("LoadWithDefaults() error = nil, want error")
		}
	})

	t.Run("invalid duration", func(t *testing.T) {
		configPath := filepath.Join(tmpDir, "invalid_duration.yaml")
		content := `
build:
  cmd: "go build"
  bin: "./main"
  delay: "invalid"
watch:
  extensions:
    - ".go"
  dirs:
    - "."
log:
  level: "info"
`
		if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
			t.Fatalf("write config file: %v", err)
		}

		_, err := LoadWithDefaults(configPath)
		if err == nil {
			t.Error("LoadWithDefaults() error = nil, want error for invalid duration")
		}
	})
}

func TestExists(t *testing.T) {
	tmpDir := t.TempDir()

	existingFile := filepath.Join(tmpDir, "existing.yaml")
	if err := os.WriteFile(existingFile, []byte("test"), 0644); err != nil {
		t.Fatalf("create file: %v", err)
	}

	if !Exists(existingFile) {
		t.Error("Exists() = false for existing file")
	}

	if Exists(filepath.Join(tmpDir, "nonexistent.yaml")) {
		t.Error("Exists() = true for nonexistent file")
	}
}

func TestWriteDefault(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "goreload.yaml")

	if err := WriteDefault(configPath); err != nil {
		t.Fatalf("WriteDefault() error = %v", err)
	}

	if !Exists(configPath) {
		t.Error("WriteDefault() did not create file")
	}

	// Verify the file can be loaded
	cfg, err := LoadWithDefaults(configPath)
	if err != nil {
		t.Errorf("LoadWithDefaults() on generated file error = %v", err)
	}

	if cfg.Build.Cmd == "" {
		t.Error("Generated config has empty build command")
	}
}
