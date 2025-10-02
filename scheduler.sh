#!/bin/bash
# Run scheduler (enqueue a job)

# $1 = payload (default jika tidak diisi)
# $2 = delay (detik, default 3)
MESSAGE=${1:-"Hello from scheduler failed"}
DELAY=${2:-3}

go run cmd/scheduler/main.go "$MESSAGE" "$DELAY"
