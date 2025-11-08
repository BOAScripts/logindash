#!/usr/bin/env bash
set -euo pipefail

# ---------------------------------------------
# LoginDash automated installer
# ---------------------------------------------
# This script downloads the latest GitHub release, extracts the
# `logindash` binary and installs it to /usr/local/bin.  It also
# creates a default configuration file in $HOME/.config/logindash
# unless one already exists.
# ---------------------------------------------

# --- Configuration --------------------------------------------
# Update these values to match your GitHub repository.
GH_USER="your-username"     # <-- replace with the GitHub username
GH_REPO="logindash"

# Asset that will be downloaded.  The release must contain a
# tar.gz with the binary named `logindash` and a config template.
ASSET_NAME="logindash-linux-amd64.tar.gz"

# --------------------------------------------
# Helper functions -------------------------------------------
# Check if a command exists.
command_exists() {
  command -v "$1" >/dev/null 2>&1
}

# --------------------------------------------
# Verify prerequisites ---------------------------------------
for cmd in curl tar; do
  if ! command_exists "$cmd"; then
    echo "Error: required tool '$cmd' is not installed." >&2
    exit 1
  fi
done

# --------------------------------------------
# Get the latest release tag -------------------------------------
RELEASE_API="https://api.github.com/repos/${GH_USER}/${GH_REPO}/releases/latest"

TAG=$(curl -s "$RELEASE_API" | grep -m1 tag_name | awk -F\" '{print $4}')
if [[ -z "$TAG" ]]; then
  echo "Error: could not determine the latest release tag." >&2
  exit 1
fi

# Construct asset download URL
ASSET_URL="https://github.com/${GH_USER}/${GH_REPO}/releases/download/${TAG}/${ASSET_NAME}"

# --------------------------------------------
# Download and extract --------------------------------------
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

echo "Downloading $ASSET_NAME (tag $TAG)..."
curl -L -o "$TMP_DIR/$ASSET_NAME" "$ASSET_URL"

echo "Extracting archive..."
tar -xzf "$TMP_DIR/$ASSET_NAME" -C "$TMP_DIR"

# --------------------------------------------
# Install binary ---------------------------------------------
BIN_PATH="/usr/local/bin/logindash"
echo "Installing binary to $BIN_PATH..."
sudo cp "$TMP_DIR/logindash" "$BIN_PATH"
sudo chmod 755 "$BIN_PATH"

# --------------------------------------------
# Install configuration -----------------------------------------
DOT_CONFIG_DIR="$HOME/.config"
CONFIG_DIR="$HOME/.config/logindash"
CONFIG_FILE="$CONFIG_DIR/config.toml"

mkdir -p "$DOT_CONFIG_DIR"
mkdir -p "$CONFIG_DIR"

if [[ ! -f "$CONFIG_FILE" ]]; then
  echo "Creating default configuration at $CONFIG_FILE..."
  cat > "$CONFIG_FILE" <<'EOF'
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
EOF
else
  echo "Config file already exists at $CONFIG_FILE. Skipping."
fi

echo
echo "Installation complete. You can now run 'logindash' from anywhere."
