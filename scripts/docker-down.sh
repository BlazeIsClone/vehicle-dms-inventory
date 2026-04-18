#!/usr/bin/env bash
set -euo pipefail

if docker compose down 2>/dev/null; then
    :
else
    echo "Falling back to Docker Compose V1"
    docker-compose down
fi
