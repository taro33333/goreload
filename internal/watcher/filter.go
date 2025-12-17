package watcher

import (
	"path/filepath"
	"strings"
)

// Filter determines whether a file path should be processed.
type Filter interface {
	// Match returns true if the path should be watched/processed.
	Match(path string) bool
}

// FilterConfig holds filter configuration.
type FilterConfig struct {
	Extensions   []string
	ExcludeDirs  []string
	ExcludeFiles []string
	Root         string
}

type filter struct {
	extensions   map[string]bool
	excludeDirs  map[string]bool
	excludeFiles []string
	root         string
}

// NewFilter creates a new Filter with the given configuration.
func NewFilter(cfg FilterConfig) Filter {
	f := &filter{
		extensions:   make(map[string]bool),
		excludeDirs:  make(map[string]bool),
		excludeFiles: cfg.ExcludeFiles,
		root:         cfg.Root,
	}

	for _, ext := range cfg.Extensions {
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		f.extensions[ext] = true
	}

	for _, dir := range cfg.ExcludeDirs {
		f.excludeDirs[dir] = true
	}

	return f
}

func (f *filter) Match(path string) bool {
	// Make path relative to root if possible.
	relPath := path
	if f.root != "" {
		if rel, err := filepath.Rel(f.root, path); err == nil {
			relPath = rel
		}
	}

	// Check excluded directories.
	if f.isInExcludedDir(relPath) {
		return false
	}

	// Check excluded files.
	if f.matchesExcludedFile(relPath) {
		return false
	}

	// Check extension.
	ext := filepath.Ext(path)
	if len(f.extensions) > 0 && !f.extensions[ext] {
		return false
	}

	return true
}

func (f *filter) isInExcludedDir(relPath string) bool {
	parts := strings.Split(filepath.ToSlash(relPath), "/")
	for _, part := range parts {
		if f.excludeDirs[part] {
			return true
		}
	}
	return false
}

func (f *filter) matchesExcludedFile(relPath string) bool {
	base := filepath.Base(relPath)
	for _, pattern := range f.excludeFiles {
		matched, err := filepath.Match(pattern, base)
		if err == nil && matched {
			return true
		}
		// Also try matching against the relative path.
		matched, err = filepath.Match(pattern, relPath)
		if err == nil && matched {
			return true
		}
	}
	return false
}

