package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/brooqs/xfce4-night-mode/internal/config"
	"github.com/brooqs/xfce4-night-mode/internal/daemon"
)

var (
	version = "dev"
)

func main() {
	var (
		daemonMode bool
		check      bool
		apply      bool
		initConfig bool
		configPath string
		showVer    bool
	)

	flag.BoolVar(&daemonMode, "daemon", false, "Run in daemon mode (continuous monitoring)")
	flag.BoolVar(&daemonMode, "d", false, "Run in daemon mode (shorthand)")
	flag.BoolVar(&check, "check", false, "Show current status (sunrise/sunset, active theme)")
	flag.BoolVar(&apply, "apply", false, "Apply the correct theme once and exit")
	flag.BoolVar(&initConfig, "init", false, "Create default config file")
	flag.StringVar(&configPath, "config", "", "Path to config file (default: ~/.config/xfce4-night-mode/config.yaml)")
	flag.BoolVar(&showVer, "version", false, "Show version")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "xfce4-night-mode — Automatic theme switching based on sunrise/sunset\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  xfce4-night-mode [flags]\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  xfce4-night-mode --init          Create default config\n")
		fmt.Fprintf(os.Stderr, "  xfce4-night-mode --check         Show current status\n")
		fmt.Fprintf(os.Stderr, "  xfce4-night-mode --apply         Apply theme once\n")
		fmt.Fprintf(os.Stderr, "  xfce4-night-mode --daemon        Run continuously\n")
	}

	flag.Parse()

	if showVer {
		fmt.Printf("xfce4-night-mode %s\n", version)
		return
	}

	// Handle --init before loading config
	if initConfig {
		path, err := config.Init()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Config file created: %s\n", path)
		fmt.Println("Edit the file to set your location and preferred themes.")
		return
	}

	// Load config
	var cfg *config.Config
	var err error

	if configPath != "" {
		cfg, err = config.LoadFrom(configPath)
	} else {
		cfg, err = config.Load()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		fmt.Fprintf(os.Stderr, "Run 'xfce4-night-mode --init' to create a default config file.\n")
		os.Exit(1)
	}

	// Execute the requested action
	switch {
	case check:
		daemon.CheckStatus(cfg)
	case apply:
		if err := daemon.ApplyOnce(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case daemonMode:
		daemon.Run(cfg)
	default:
		flag.Usage()
		os.Exit(1)
	}
}
