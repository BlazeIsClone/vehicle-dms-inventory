#!/usr/bin/env bash
set -euo pipefail

if docker compose up --build 2>/dev/null; then
    :
else
    echo "Falling back to Docker Compose V1"
    docker-compose up --build
fi
