# xfce4-night-mode

Automatic theme switching for XFCE4 based on sunrise and sunset times.

The daemon calculates sunrise and sunset for your location, then seamlessly switches between your preferred light and dark themes. No cron jobs needed — it runs as a lightweight background service and sleeps until the next transition.

## Features

- 🌅 **Sunrise/sunset calculation** based on GPS coordinates
- 🎨 **Switches GTK, icon, and window manager themes** via `xfconf-query`
- 🔄 **Daemon mode** — runs continuously, sleeps until next transition
- ⚡ **Single-shot mode** — apply the correct theme once and exit
- 📋 **Status check** — see current mode, sunrise/sunset times, and active theme
- 🛠️ **Systemd integration** — start automatically on login

## Installation

### Prerequisites

- XFCE4 desktop environment
- Go 1.21+ (for building from source)
- `xfconf-query` (comes pre-installed with XFCE4)

### Build & Install

```bash
git clone https://github.com/brooqs/xfce4-night-mode.git
cd xfce4-night-mode
make install
```

This installs the binary to `~/.local/bin/` and the systemd service file.

## Configuration

### Initialize

```bash
xfce4-night-mode --init
```

This creates `~/.config/xfce4-night-mode/config.yaml` with default values.

### Edit Config

```yaml
# Your GPS coordinates (default: Istanbul)
location:
  latitude: 41.0082
  longitude: 28.9784

# Theme applied during daytime
day_theme:
  gtk_theme: "Adwaita"
  icon_theme: "Adwaita"
  wm_theme: "Default"

# Theme applied at night
night_theme:
  gtk_theme: "Adwaita-dark"
  icon_theme: "Adwaita"
  wm_theme: "Default-hdpi"

# How often to check (minutes). The daemon also wakes at transitions.
check_interval: 5
```

> **Tip:** Find your coordinates at [latlong.net](https://www.latlong.net)

## Usage

```bash
# Show current status (sunrise/sunset times, active theme)
xfce4-night-mode --check

# Apply the correct theme once and exit
xfce4-night-mode --apply

# Run as a daemon (continuous monitoring)
xfce4-night-mode --daemon

# Use a custom config file
xfce4-night-mode --config /path/to/config.yaml --daemon
```

### Example Output

```
$ xfce4-night-mode --check
Current mode:    Night 🌙
Sunrise today:   06:48:12
Sunset today:    17:55:30
Next Sunrise:    06:47:05 (in 6h32m)

Night theme: GTK=Adwaita-dark  Icon=Adwaita  WM=Default-hdpi
Current:     GTK=Adwaita-dark  Icon=Adwaita  WM=Default-hdpi
```

## Autostart with Systemd

```bash
# Enable and start the service
make enable

# Check service status
systemctl --user status xfce4-night-mode

# View logs
journalctl --user -u xfce4-night-mode -f

# Disable the service
make disable
```

## Uninstall

```bash
make uninstall
```

## How It Works

1. Reads your location from the config file
2. Calculates today's sunrise and sunset times using the [go-sunrise](https://github.com/nathan-osman/go-sunrise) library
3. Determines if it's currently daytime or nighttime
4. Applies the appropriate theme using `xfconf-query`
5. In daemon mode, sleeps until the next transition and repeats

## License

MIT
