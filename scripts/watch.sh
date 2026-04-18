#!/usr/bin/env bash
set -euo pipefail

if command -v air > /dev/null; then
    echo "Watching..."
    air
else
    read -rp "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice
    if [ "$choice" != "n" ] && [ "$choice" != "N" ]; then
        go install github.com/air-verse/air@latest
        echo "Watching..."
        air
    else
        echo "You chose not to install air. Exiting..."
        exit 1
    fi
fi
