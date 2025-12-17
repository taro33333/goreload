package watcher

import (
	"path/filepath"
	"testing"
)

func TestNewFilter(t *testing.T) {
	cfg := FilterConfig{
		Extensions:   []string{".go", "html"},
		ExcludeDirs:  []string{"vendor", "tmp"},
		ExcludeFiles: []string{"*_test.go"},
		Root:         "/project",
	}

	f := NewFilter(cfg)
	if f == nil {
		t.Error("NewFilter() returned nil")
	}
}

func TestFilter_Match(t *testing.T) {
	f := NewFilter(FilterConfig{
		Extensions:   []string{".go"},
		ExcludeDirs:  []string{"vendor", "tmp", ".git"},
		ExcludeFiles: []string{"*_test.go", "*.gen.go"},
		Root:         "/project",
	})

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "go file in root",
			path: "/project/main.go",
			want: true,
		},
		{
			name: "go file in subdirectory",
			path: "/project/pkg/handler.go",
			want: true,
		},
		{
			name: "test file excluded",
			path: "/project/main_test.go",
			want: false,
		},
		{
			name: "generated file excluded",
			path: "/project/types.gen.go",
			want: false,
		},
		{
			name: "vendor directory excluded",
			path: "/project/vendor/lib/lib.go",
			want: false,
		},
		{
			name: "tmp directory excluded",
			path: "/project/tmp/main.go",
			want: false,
		},
		{
			name: "git directory excluded",
			path: "/project/.git/config",
			want: false,
		},
		{
			name: "non-go file excluded",
			path: "/project/README.md",
			want: false,
		},
		{
			name: "yaml file excluded",
			path: "/project/config.yaml",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := f.Match(tt.path); got != tt.want {
				t.Errorf("Match(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestFilter_Match_MultipleExtensions(t *testing.T) {
	f := NewFilter(FilterConfig{
		Extensions:   []string{".go", ".html", ".tmpl"},
		ExcludeDirs:  []string{},
		ExcludeFiles: []string{},
		Root:         "/project",
	})

	tests := []struct {
		path string
		want bool
	}{
		{"/project/main.go", true},
		{"/project/index.html", true},
		{"/project/base.tmpl", true},
		{"/project/style.css", false},
		{"/project/script.js", false},
	}

	for _, tt := range tests {
		t.Run(filepath.Base(tt.path), func(t *testing.T) {
			if got := f.Match(tt.path); got != tt.want {
				t.Errorf("Match(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestFilter_Match_ExtensionNormalization(t *testing.T) {
	// Extensions without leading dot should be normalized
	f := NewFilter(FilterConfig{
		Extensions:   []string{"go", ".html"},
		ExcludeDirs:  []string{},
		ExcludeFiles: []string{},
		Root:         "/project",
	})

	tests := []struct {
		path string
		want bool
	}{
		{"/project/main.go", true},
		{"/project/index.html", true},
	}

	for _, tt := range tests {
		t.Run(filepath.Base(tt.path), func(t *testing.T) {
			if got := f.Match(tt.path); got != tt.want {
				t.Errorf("Match(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestFilter_Match_NestedExcludeDirs(t *testing.T) {
	f := NewFilter(FilterConfig{
		Extensions:   []string{".go"},
		ExcludeDirs:  []string{"node_modules"},
		ExcludeFiles: []string{},
		Root:         "/project",
	})

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "deeply nested node_modules",
			path: "/project/frontend/node_modules/pkg/index.go",
			want: false,
		},
		{
			name: "regular nested directory",
			path: "/project/internal/handler/handler.go",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := f.Match(tt.path); got != tt.want {
				t.Errorf("Match(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestFilter_Match_EmptyRoot(t *testing.T) {
	f := NewFilter(FilterConfig{
		Extensions:   []string{".go"},
		ExcludeDirs:  []string{"vendor"},
		ExcludeFiles: []string{},
		Root:         "",
	})

	// Should still work with absolute paths
	if !f.Match("/some/path/main.go") {
		t.Error("Match() should work with empty root")
	}
}

func TestFilter_Match_GlobPatterns(t *testing.T) {
	f := NewFilter(FilterConfig{
		Extensions:   []string{".go"},
		ExcludeDirs:  []string{},
		ExcludeFiles: []string{"mock_*.go", "*.pb.go"},
		Root:         "/project",
	})

	tests := []struct {
		path string
		want bool
	}{
		{"/project/mock_service.go", false},
		{"/project/user.pb.go", false},
		{"/project/service.go", true},
		{"/project/handler.go", true},
	}

	for _, tt := range tests {
		t.Run(filepath.Base(tt.path), func(t *testing.T) {
			if got := f.Match(tt.path); got != tt.want {
				t.Errorf("Match(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}
