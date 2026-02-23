package solar

import (
	"fmt"
	"time"

	"github.com/nathan-osman/go-sunrise"
)

// Info holds sunrise/sunset times and daytime status.
type Info struct {
	Sunrise   time.Time
	Sunset    time.Time
	IsDaytime bool
}

// Calculate returns sunrise/sunset info for the given location and current time.
func Calculate(lat, lon float64) Info {
	now := time.Now()
	rise, set := sunrise.SunriseSunset(lat, lon, now.Year(), now.Month(), now.Day())

	// Convert to local timezone
	loc := now.Location()
	rise = rise.In(loc)
	set = set.In(loc)

	isDaytime := now.After(rise) && now.Before(set)

	return Info{
		Sunrise:   rise,
		Sunset:    set,
		IsDaytime: isDaytime,
	}
}

// NextTransition returns the time of the next sunrise or sunset.
// If it's daytime, returns the next sunset. If nighttime, returns
// the next sunrise (which might be tomorrow).
func NextTransition(lat, lon float64) time.Time {
	now := time.Now()
	info := Calculate(lat, lon)

	if info.IsDaytime {
		// Next transition is sunset today
		return info.Sunset
	}

	// It's nighttime — if we're past sunset, next transition is tomorrow's sunrise
	if now.After(info.Sunset) || now.Equal(info.Sunset) {
		tomorrow := now.AddDate(0, 0, 1)
		rise, _ := sunrise.SunriseSunset(lat, lon, tomorrow.Year(), tomorrow.Month(), tomorrow.Day())
		return rise.In(now.Location())
	}

	// Before sunrise today
	return info.Sunrise
}

// FormatInfo returns a human-readable string of the solar info.
func FormatInfo(lat, lon float64) string {
	info := Calculate(lat, lon)
	next := NextTransition(lat, lon)

	mode := "Night 🌙"
	if info.IsDaytime {
		mode = "Day ☀️"
	}

	nextLabel := "Sunset"
	if !info.IsDaytime {
		nextLabel = "Sunrise"
	}

	return fmt.Sprintf(
		"Current mode:    %s\n"+
			"Sunrise today:   %s\n"+
			"Sunset today:    %s\n"+
			"Next %s: %s (in %s)",
		mode,
		info.Sunrise.Format("15:04:05"),
		info.Sunset.Format("15:04:05"),
		nextLabel,
		next.Format("15:04:05"),
		time.Until(next).Round(time.Second),
	)
}
