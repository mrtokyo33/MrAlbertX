#!/bin/bash
set -e

echo "▶Compiling Mr. Albert X for Linux..."

go build -o mr-x .

echo "Compilation successful."
echo "▶Installing mr-x command..."

INSTALL_DIR="$HOME/.local/bin"
mkdir -p "$INSTALL_DIR"

mv mr-x "$INSTALL_DIR/"

echo "Mr. Albert X installed successfully!"
echo
echo "You can now run 'mr-x' from anywhere in your terminal."
echo "If the command is not found, please restart your terminal session for the PATH changes to take effect."
