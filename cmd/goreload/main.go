// Package main provides the CLI entry point for goreload.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/user/goreload/internal/config"
	"github.com/user/goreload/internal/engine"
	"github.com/user/goreload/internal/logger"
)

// Version is set at build time.
var Version = "v0.1.0"

func main() {
	if err := rootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func rootCmd() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:   "goreload",
		Short: "Hot reload for Go applications",
		Long:  "goreload watches your Go files and automatically rebuilds and restarts your application.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(configPath)
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", config.DefaultConfigFile, "config file path")

	cmd.AddCommand(versionCmd())
	cmd.AddCommand(initCmd())

	return cmd
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("goreload %s\n", Version)
		},
	}
}

func initCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Generate a default configuration file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if config.Exists(config.DefaultConfigFile) {
				return fmt.Errorf("%s already exists", config.DefaultConfigFile)
			}

			if err := config.WriteDefault(config.DefaultConfigFile); err != nil {
				return fmt.Errorf("write config: %w", err)
			}

			fmt.Printf("Created %s\n", config.DefaultConfigFile)
			return nil
		},
	}
}

func run(configPath string) error {
	// Load configuration.
	var cfg *config.Config
	var err error

	if config.Exists(configPath) {
		cfg, err = config.LoadWithDefaults(configPath)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
	} else if configPath != config.DefaultConfigFile {
		// User specified a non-default config file that doesn't exist.
		return fmt.Errorf("config file not found: %s", configPath)
	} else {
		// Use defaults if no config file exists.
		cfg = config.Default()
	}

	// Create logger.
	log := logger.New(logger.Config{
		Color: cfg.Log.Color,
		Time:  cfg.Log.Time,
		Level: cfg.Log.Level,
	})

	// Print banner.
	logger.Banner(os.Stdout, Version)

	// Create engine.
	eng, err := engine.New(cfg, log)
	if err != nil {
		return fmt.Errorf("create engine: %w", err)
	}

	// Setup signal handling.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		cancel()
	}()

	// Run engine.
	if err := eng.Run(ctx); err != nil && err != context.Canceled {
		return err
	}

	return nil
}
