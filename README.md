# LoginDash

A lightweight terminal dashboard that shows real‑time system information when you log in via SSH.
It is written in Go and renders a coloured, easy‑to‑read summary of:

- **System** – uptime, CPU usage, RAM usage
- **Network** – IP address, default gateway, DNS servers
- **Storage** – root and mounted drives
- **Services** – status of user‑selected systemd services

The dashboard is configurable via a TOML file and can be integrated into your
`~/.ssh/rc` or `~/.bashrc` so that it displays automatically on login.

## Features

| Feature | Description |
|---------|-------------|
| **Dynamic storage** | Scan `/mnt` for new mounts automatically. |
| **Service monitoring** | Show the status of any systemd service you care about. |
| **Portable** | Standard Linux utilities depedencies (`top`, `free`, `df`, `systemctl`, etc.). |
| **Custom labels** | Set the width of the labels and colour thresholds in the config. |

## Installation

### Automated

> Review the `install.sh` code

```bash
curl -fsSL "https://raw.githubusercontent.com/BOAScripts/logindash/refs/heads/main/install/install.sh" | bash
```

### Build it yourself

1. **Clone the repo**

   ```bash
   git clone https://github.com/BOAScripts/logindash.git
   cd logindash/
   ```

2. **Build**

   ```bash
   go build -o logindash
   ```

3. **Place the binary**

   ```bash
   sudo mv logindash /usr/local/bin/
   ```

4. **Copy/Edit the config file**

  ```bash
  mkdir -p ~/.config/logindash
  cp config/config.toml .config/logindash/config.toml
  ```

## Usage

```bash
logindash
```

### Options

| Options | Description |
| -- | -- |
| `--config <string>` | Path to the configuration file|
| `-h`, `--help` | Show help |

Default config path: `~/.config/logindash/config.toml`.

```bash
logindash -config ~/.mylogindashconfig.toml
```

### Auto‑run on SSH login

Add the following line at t he end of your `~/.bashrc`:

```bash
if [ -f logindash ]; then
    logindash
fi
```

## Configuration

Open `~/.config/logindash/config.toml`:

```toml
[display]
label_width  = 15 # <- Default value if ommitted
green_until  = 65 # <- Default value if ommitted
orange_until = 85 # <- Default value if ommitted

[disks]
paths = [
  "/home",
  "/var"
]

[services]
monitored = [
  "ssh",
  "cron"
]
```

- **display.label_width** – Width of the left‑hand label column. This ensure the values to be on the same x-axis
- **display.green_until / orange_until** – Thresholds for colour coding the usage bars. (green from 0 to 65, orange from 66 to 85, rest is red)
- **disks.paths** – Additional paths to display disk usage for.
- **services.monitored** – Systemd services to monitor.

## Customisation

- **Thresholds** – Adjust `green_until` and `orange_until` to suit your monitoring style.
- **Mount detection** – The program automatically scans `/mnt` for new mounts. You can add paths in the config as needed.
- **Colours** – Edit the `lipgloss` styles in `main.go` to change the palette.

## Disclaimer

> This has been written with the help of some AIs (Claude Sonnet 4.5 & gpt-oss). Sorry for the slop
