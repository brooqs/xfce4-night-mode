package daemon

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/brooqs/xfce4-night-mode/internal/config"
	"github.com/brooqs/xfce4-night-mode/internal/solar"
	"github.com/brooqs/xfce4-night-mode/internal/xfce"
)

// Run starts the daemon loop. It checks the current time against
// sunrise/sunset and applies the appropriate theme, then sleeps
// until the next transition.
func Run(cfg *config.Config) {
	log.SetFlags(log.LstdFlags | log.Lmsgprefix)
	log.SetPrefix("[xfce4-night-mode] ")

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Daemon started")
	log.Printf("Location: %.4f, %.4f", cfg.Location.Latitude, cfg.Location.Longitude)

	var lastMode string

	for {
		mode, err := applyCorrectTheme(cfg)
		if err != nil {
			log.Printf("Error applying theme: %v", err)
		} else if mode != lastMode {
			log.Printf("Theme switched to: %s", mode)
			lastMode = mode
		}

		// Calculate sleep duration
		sleepDuration := calculateSleep(cfg)
		log.Printf("Next check in %s", sleepDuration.Round(time.Second))

		timer := time.NewTimer(sleepDuration)
		select {
		case <-timer.C:
			// Time to check again
		case sig := <-sigCh:
			timer.Stop()
			log.Printf("Received signal %v, shutting down", sig)
			return
		}
	}
}

// ApplyOnce checks the current time and applies the correct theme once.
func ApplyOnce(cfg *config.Config) error {
	mode, err := applyCorrectTheme(cfg)
	if err != nil {
		return err
	}
	fmt.Printf("Applied %s theme\n", mode)
	return nil
}

// CheckStatus prints the current status without applying any changes.
func CheckStatus(cfg *config.Config) {
	lat := cfg.Location.Latitude
	lon := cfg.Location.Longitude

	fmt.Println(solar.FormatInfo(lat, lon))
	fmt.Println()

	info := solar.Calculate(lat, lon)
	if info.IsDaytime {
		fmt.Printf("Day theme:   GTK=%s  Icon=%s  WM=%s\n",
			cfg.DayTheme.GtkTheme, cfg.DayTheme.IconTheme, cfg.DayTheme.WmTheme)
	} else {
		fmt.Printf("Night theme: GTK=%s  Icon=%s  WM=%s\n",
			cfg.NightTheme.GtkTheme, cfg.NightTheme.IconTheme, cfg.NightTheme.WmTheme)
	}

	// Try to show current theme
	current, err := xfce.GetCurrentTheme()
	if err == nil {
		fmt.Printf("\nCurrent:     GTK=%s  Icon=%s  WM=%s\n",
			current.GtkTheme, current.IconTheme, current.WmTheme)
	}
}

// applyCorrectTheme determines the current mode and applies the matching theme.
func applyCorrectTheme(cfg *config.Config) (string, error) {
	info := solar.Calculate(cfg.Location.Latitude, cfg.Location.Longitude)

	var theme config.ThemeConfig
	var mode string

	if info.IsDaytime {
		theme = cfg.DayTheme
		mode = "day"
	} else {
		theme = cfg.NightTheme
		mode = "night"
	}

	err := xfce.ApplyTheme(xfce.ThemeConfig{
		GtkTheme:  theme.GtkTheme,
		IconTheme: theme.IconTheme,
		WmTheme:   theme.WmTheme,
	})

	return mode, err
}

// calculateSleep returns how long to sleep before the next check.
// It sleeps until the next sunrise/sunset transition, but caps at
// the configured check interval.
func calculateSleep(cfg *config.Config) time.Duration {
	next := solar.NextTransition(cfg.Location.Latitude, cfg.Location.Longitude)
	untilTransition := time.Until(next)

	maxSleep := time.Duration(cfg.CheckInterval) * time.Minute

	// If the transition is soon, sleep until just past it
	if untilTransition > 0 && untilTransition < maxSleep {
		return untilTransition + 10*time.Second // small buffer
	}

	return maxSleep
}
