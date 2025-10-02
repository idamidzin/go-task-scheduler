# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A Redis-backed distributed task scheduler written in Go. The system uses a **Scheduler-Worker pattern** where schedulers enqueue delayed jobs into Redis sorted sets, and workers poll and execute jobs when they become ready.

## Architecture

### Core Components

- **Scheduler** ([internal/scheduler/scheduler.go](internal/scheduler/scheduler.go)): Enqueues jobs into Redis sorted set with delay-based scoring
- **Worker** ([internal/worker/worker.go](internal/worker/worker.go)): Polls Redis for ready jobs and executes them with retry logic
- **Job** ([internal/job/job.go](internal/job/job.go)): Job structure with ID, Payload, Delay, Retries, and MaxRetry
- **Config** ([internal/config/config.go](internal/config/config.go)): Environment-based configuration loader

### Job Flow

1. Scheduler receives job with payload and delay
2. Job is assigned UUID and pushed to Redis sorted set "jobs" with score = `now + delay` (in milliseconds)
3. Worker polls Redis every `INTERVAL` for jobs with score ≤ current time
4. Worker removes job from queue and executes it
5. If payload equals "fail", job is retried up to `MaxRetry` times with 2-second delay
6. Successful jobs are logged and removed; exhausted retries are discarded

### Redis Usage

- **Key**: `jobs` (sorted set)
- **Score**: Unix millisecond timestamp when job becomes ready
- **Member**: JSON-encoded job struct

## Development Commands

### Running the System

```bash
# Start worker (in one terminal)
./worker.sh
# or
go run cmd/worker/main.go

# Enqueue a job (in another terminal)
./scheduler.sh "my payload" 5
# or
go run cmd/scheduler/main.go "my payload" 5
```

### Building

```bash
# Build both binaries
go build -o bin/scheduler cmd/scheduler/main.go
go build -o bin/worker cmd/worker/main.go

# Run built binaries
./bin/worker
./bin/scheduler "test job" 3
```

### Testing

```bash
# Run all tests
go test ./...

# Test specific package
go test ./internal/scheduler
go test ./internal/worker

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...
```

### Dependencies

```bash
# Install dependencies
go mod download

# Update dependencies
go mod tidy

# Verify dependencies
go mod verify
```

## Configuration

Environment variables (create `.env` file):

```bash
REDIS_ADDR=localhost:6379  # Redis server address
REDIS_DB=0                 # Redis database number
INTERVAL=200ms             # Worker poll interval
JOB_DELAY=3s              # Default job delay when scheduling
JOB_MAX_RETRY=3           # Maximum retry attempts
```

## Scheduler CLI Arguments

```bash
go run cmd/scheduler/main.go [payload] [delay_in_seconds]

# Examples:
go run cmd/scheduler/main.go "success job" 5
go run cmd/scheduler/main.go "fail" 2  # Will trigger retry logic
```

## Testing Job Behavior

- **Success case**: Any payload except "fail"
- **Failure case**: Payload = "fail" triggers retry logic (retries up to `MaxRetry` with 2s delay)
- Worker logs show job execution status with emoji indicators (✅ success, ⚠️ retry, ❌ error)

## Module Information

- **Module path**: `github.com/idamidzin/go-task-scheduler`
- **Go version**: 1.19
- **Key dependencies**:
  - `github.com/redis/go-redis/v9` - Redis client
  - `github.com/google/uuid` - UUID generation
  - `github.com/joho/godotenv` - Environment variable loading
