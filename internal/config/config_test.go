package config

import (
	"testing"
	"time"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr error
	}{
		{
			name:    "valid config",
			cfg:     *validConfig(),
			wantErr: nil,
		},
		{
			name: "empty build cmd",
			cfg: func() Config {
				c := *validConfig()
				c.Build.Cmd = ""
				return c
			}(),
			wantErr: ErrEmptyBuildCmd,
		},
		{
			name: "empty bin",
			cfg: func() Config {
				c := *validConfig()
				c.Build.Bin = ""
				return c
			}(),
			wantErr: ErrEmptyBin,
		},
		{
			name: "negative delay",
			cfg: func() Config {
				c := *validConfig()
				c.Build.Delay = -1 * time.Second
				return c
			}(),
			wantErr: ErrInvalidDelay,
		},
		{
			name: "negative kill delay",
			cfg: func() Config {
				c := *validConfig()
				c.Build.KillDelay = -1 * time.Second
				return c
			}(),
			wantErr: ErrInvalidKillDelay,
		},
		{
			name: "no extensions",
			cfg: func() Config {
				c := *validConfig()
				c.Watch.Extensions = []string{}
				return c
			}(),
			wantErr: ErrNoExtensions,
		},
		{
			name: "no dirs",
			cfg: func() Config {
				c := *validConfig()
				c.Watch.Dirs = []string{}
				return c
			}(),
			wantErr: ErrNoDirs,
		},
		{
			name: "invalid log level",
			cfg: func() Config {
				c := *validConfig()
				c.Log.Level = "invalid"
				return c
			}(),
			wantErr: ErrInvalidLogLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("Validate() error = %v, want nil", err)
				}
				return
			}
			if err == nil {
				t.Errorf("Validate() error = nil, want %v", tt.wantErr)
				return
			}
			// Check if the error contains the expected error
			if err.Error() == "" {
				t.Errorf("Validate() error message is empty")
			}
		})
	}
}

func TestConfig_AbsRoot(t *testing.T) {
	tests := []struct {
		name    string
		root    string
		wantAbs bool
	}{
		{
			name:    "empty root returns cwd",
			root:    "",
			wantAbs: true,
		},
		{
			name:    "dot returns cwd",
			root:    ".",
			wantAbs: true,
		},
		{
			name:    "relative path",
			root:    "testdir",
			wantAbs: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{Root: tt.root}
			got, err := cfg.AbsRoot()
			if err != nil {
				t.Errorf("AbsRoot() error = %v", err)
				return
			}
			if tt.wantAbs && got[0] != '/' {
				t.Errorf("AbsRoot() = %v, want absolute path", got)
			}
		})
	}
}

func TestLogConfig_validate(t *testing.T) {
	validLevels := []string{"debug", "info", "warn", "error"}
	for _, level := range validLevels {
		t.Run("valid_"+level, func(t *testing.T) {
			cfg := LogConfig{Level: level}
			if err := cfg.validate(); err != nil {
				t.Errorf("validate() error = %v for level %s", err, level)
			}
		})
	}

	t.Run("invalid level", func(t *testing.T) {
		cfg := LogConfig{Level: "trace"}
		if err := cfg.validate(); err != ErrInvalidLogLevel {
			t.Errorf("validate() error = %v, want %v", err, ErrInvalidLogLevel)
		}
	})
}

func validConfig() *Config {
	return &Config{
		Root:   ".",
		TmpDir: "tmp",
		Build: BuildConfig{
			Cmd:       "go build -o ./tmp/main .",
			Bin:       "./tmp/main",
			Args:      []string{},
			Delay:     200 * time.Millisecond,
			KillDelay: 500 * time.Millisecond,
		},
		Watch: WatchConfig{
			Extensions:   []string{".go"},
			Dirs:         []string{"."},
			ExcludeDirs:  []string{"tmp"},
			ExcludeFiles: []string{},
		},
		Log: LogConfig{
			Color: true,
			Time:  true,
			Level: "info",
		},
	}
}
