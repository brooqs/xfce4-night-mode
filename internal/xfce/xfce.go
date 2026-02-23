package xfce

import (
	"fmt"
	"os/exec"
	"strings"
)

// ThemeConfig mirrors config.ThemeConfig for applying themes.
type ThemeConfig struct {
	GtkTheme  string
	IconTheme string
	WmTheme   string
}

// themeProperty represents a single xfconf-query target.
type themeProperty struct {
	channel  string
	property string
	value    string
	label    string
}

// ApplyTheme applies the given theme configuration using xfconf-query.
func ApplyTheme(theme ThemeConfig) error {
	props := []themeProperty{
		{"xsettings", "/Net/ThemeName", theme.GtkTheme, "GTK Theme"},
		{"xsettings", "/Net/IconThemeName", theme.IconTheme, "Icon Theme"},
		{"xfwm4", "/general/theme", theme.WmTheme, "WM Theme"},
	}

	var errors []string
	for _, p := range props {
		if p.value == "" {
			continue
		}
		if err := xfconfSet(p.channel, p.property, p.value); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", p.label, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("some themes failed to apply:\n  %s", strings.Join(errors, "\n  "))
	}
	return nil
}

// GetCurrentTheme reads the current XFCE4 theme settings.
func GetCurrentTheme() (*ThemeConfig, error) {
	gtk, err := xfconfGet("xsettings", "/Net/ThemeName")
	if err != nil {
		return nil, fmt.Errorf("could not read GTK theme: %w", err)
	}

	icon, err := xfconfGet("xsettings", "/Net/IconThemeName")
	if err != nil {
		return nil, fmt.Errorf("could not read icon theme: %w", err)
	}

	wm, err := xfconfGet("xfwm4", "/general/theme")
	if err != nil {
		return nil, fmt.Errorf("could not read WM theme: %w", err)
	}

	return &ThemeConfig{
		GtkTheme:  gtk,
		IconTheme: icon,
		WmTheme:   wm,
	}, nil
}

// xfconfSet sets a property value via xfconf-query.
func xfconfSet(channel, property, value string) error {
	cmd := exec.Command("xfconf-query", "-c", channel, "-p", property, "-s", value)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s (output: %s)", err, strings.TrimSpace(string(output)))
	}
	return nil
}

// xfconfGet reads a property value via xfconf-query.
func xfconfGet(channel, property string) (string, error) {
	cmd := exec.Command("xfconf-query", "-c", channel, "-p", property)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s (output: %s)", err, strings.TrimSpace(string(output)))
	}
	return strings.TrimSpace(string(output)), nil
}
