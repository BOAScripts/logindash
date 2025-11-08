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
curl -fsSL ...| bash
```

### Build it yourself

1. **Clone the repo**

   ```bash
   git clone https://github.com/yourusername/logindash.git
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

- `--config string` – Path to the configuration file.
  Default: `~/.config/logindash/config.toml`.
- `-h`, `--help` – Show help.

```bash
logindash -config ~/.config/logindash/config.toml
```

### Auto‑run on SSH login

Add the following line to `~/.ssh/rc` (or `~/.bashrc` if you use that):

```bash
/usr/local/bin/logindash
```

## Configuration

Create the directory `~/.config/logindash/` and add `config.toml`:

```toml
[display]
label_width  = 15
green_until  = 65
orange_until = 85

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

- **display.label_width** – Width of the left‑hand label column.
- **display.green_until / orange_until** – Thresholds for colour coding the usage bars.
- **disks.paths** – Additional paths to display disk usage for.
- **services.monitored** – Systemd services to monitor.

## Customisation

- **Colours** – Edit the `lipgloss` styles in `main.go` to change the palette.
- **Thresholds** – Adjust `green_until` and `orange_until` to suit your monitoring style.
- **Mount detection** – The program automatically scans `/mnt` for new mounts; you can add or remove paths in the config as needed.

## Disclaimer

> This has been written with the help of some AIs (Claude Sonnet 4.5 & gpt-oss). Sorry for the slop
