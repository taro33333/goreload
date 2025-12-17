package builder

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	cfg := Config{
		Cmd:    "go build",
		Bin:    "./main",
		TmpDir: "tmp",
		Root:   ".",
	}

	b := New(cfg)
	if b == nil {
		t.Error("New() returned nil")
	}
}

func TestBuilder_Build(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("successful build", func(t *testing.T) {
		buildDir := filepath.Join(tmpDir, "buildtest")
		if err := os.MkdirAll(buildDir, 0755); err != nil {
			t.Fatalf("create build dir: %v", err)
		}

		// Create go.mod
		goMod := filepath.Join(buildDir, "go.mod")
		if err := os.WriteFile(goMod, []byte("module testapp\n\ngo 1.21\n"), 0644); err != nil {
			t.Fatalf("write go.mod: %v", err)
		}

		// Create a simple Go file
		goFile := filepath.Join(buildDir, "main.go")
		content := `package main
func main() {}
`
		if err := os.WriteFile(goFile, []byte(content), 0644); err != nil {
			t.Fatalf("write go file: %v", err)
		}

		binPath := filepath.Join(buildDir, "main")
		b := New(Config{
			Cmd:    "go build -o " + binPath + " .",
			Bin:    binPath,
			TmpDir: buildDir,
			Root:   buildDir,
		})

		ctx := context.Background()
		result := b.Build(ctx)

		if !result.Success {
			t.Errorf("Build() success = false, error = %v, output = %s", result.Error, result.Output)
		}
		if result.Duration <= 0 {
			t.Error("Build() duration should be positive")
		}

		// Check binary exists
		if _, err := os.Stat(binPath); os.IsNotExist(err) {
			t.Error("Build() did not create binary")
		}
	})

	t.Run("build failure", func(t *testing.T) {
		// Create invalid Go file
		badDir := filepath.Join(tmpDir, "bad")
		if err := os.MkdirAll(badDir, 0755); err != nil {
			t.Fatalf("create bad dir: %v", err)
		}

		// Create go.mod
		goMod := filepath.Join(badDir, "go.mod")
		if err := os.WriteFile(goMod, []byte("module badapp\n\ngo 1.21\n"), 0644); err != nil {
			t.Fatalf("write go.mod: %v", err)
		}

		goFile := filepath.Join(badDir, "main.go")
		content := `package main
func main() {
	invalid syntax here
}
`
		if err := os.WriteFile(goFile, []byte(content), 0644); err != nil {
			t.Fatalf("write go file: %v", err)
		}

		b := New(Config{
			Cmd:    "go build -o ./bad_main .",
			Bin:    "./bad_main",
			TmpDir: badDir,
			Root:   badDir,
		})

		ctx := context.Background()
		result := b.Build(ctx)

		if result.Success {
			t.Error("Build() should fail for invalid Go code")
		}
		if result.Error == nil {
			t.Error("Build() should return error for invalid Go code")
		}
	})

	t.Run("context cancellation", func(t *testing.T) {
		b := New(Config{
			Cmd:    "sleep 10",
			Bin:    "./main",
			TmpDir: tmpDir,
			Root:   tmpDir,
		})

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		result := b.Build(ctx)

		if result.Success {
			t.Error("Build() should fail on context cancellation")
		}
	})

	t.Run("empty command", func(t *testing.T) {
		b := New(Config{
			Cmd:    "",
			Bin:    "./main",
			TmpDir: tmpDir,
			Root:   tmpDir,
		})

		ctx := context.Background()
		result := b.Build(ctx)

		if result.Success {
			t.Error("Build() should fail for empty command")
		}
	})
}

func TestBuilder_Clean(t *testing.T) {
	tmpDir := t.TempDir()

	binPath := filepath.Join(tmpDir, "main")
	if err := os.WriteFile(binPath, []byte("binary"), 0755); err != nil {
		t.Fatalf("create binary: %v", err)
	}

	b := New(Config{
		Cmd:    "go build",
		Bin:    binPath,
		TmpDir: tmpDir,
		Root:   tmpDir,
	})

	if err := b.Clean(); err != nil {
		t.Errorf("Clean() error = %v", err)
	}

	if _, err := os.Stat(binPath); !os.IsNotExist(err) {
		t.Error("Clean() did not remove binary")
	}

	// Clean again should not error
	if err := b.Clean(); err != nil {
		t.Errorf("Clean() on nonexistent file error = %v", err)
	}
}

func TestParseCommand(t *testing.T) {
	tests := []struct {
		name string
		cmd  string
		want []string
	}{
		{
			name: "simple command",
			cmd:  "go build",
			want: []string{"go", "build"},
		},
		{
			name: "command with flags",
			cmd:  "go build -o ./main .",
			want: []string{"go", "build", "-o", "./main", "."},
		},
		{
			name: "quoted argument",
			cmd:  `go build -ldflags "-X main.version=1.0"`,
			want: []string{"go", "build", "-ldflags", "-X main.version=1.0"},
		},
		{
			name: "single quoted argument",
			cmd:  `go build -ldflags '-X main.version=1.0'`,
			want: []string{"go", "build", "-ldflags", "-X main.version=1.0"},
		},
		{
			name: "empty command",
			cmd:  "",
			want: []string{},
		},
		{
			name: "multiple spaces",
			cmd:  "go   build   -o   main",
			want: []string{"go", "build", "-o", "main"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseCommand(tt.cmd)
			if len(got) != len(tt.want) {
				t.Errorf("parseCommand(%q) = %v, want %v", tt.cmd, got, tt.want)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("parseCommand(%q)[%d] = %v, want %v", tt.cmd, i, got[i], tt.want[i])
				}
			}
		})
	}
}
