# Go Task Scheduler

A lightweight, Redis-backed distributed task scheduler built with Go. This system allows you to schedule delayed job execution with automatic retry mechanisms using a Scheduler-Worker architecture pattern.

## Features

- **Delayed Job Execution**: Schedule jobs to run after a specified delay
- **Automatic Retry Logic**: Failed jobs are automatically retried with configurable max retry attempts
- **Redis-Backed Queue**: Uses Redis sorted sets for reliable, distributed job storage
- **Worker Polling**: Efficient worker polling with configurable intervals
- **UUID Job Tracking**: Each job is assigned a unique identifier for tracking
- **Environment Configuration**: Flexible configuration via environment variables
- **Graceful Shutdown**: Worker supports graceful shutdown with signal handling (SIGINT, SIGTERM)
- **Job Status Logging**: Clear emoji-based logging for job execution status

## Architecture

### Components

1. **Scheduler** - Enqueues jobs into Redis with delay-based scheduling
2. **Worker** - Polls Redis for ready jobs and executes them
3. **Redis** - Acts as the job queue using sorted sets (score = execution timestamp)

### How It Works

1. Scheduler receives a job with a payload and delay
2. Job is stored in Redis sorted set with score = `current_time + delay`
3. Worker polls Redis every `INTERVAL` for jobs with score ≤ current time
4. Worker executes the job or retries on failure (payload = "fail")
5. Successful jobs are logged and removed from the queue

## Prerequisites

- Go 1.19 or higher
- Redis server running (default: `localhost:6379`)

## Installation

```bash
# Clone the repository
git clone https://github.com/idamidzin/go-task-scheduler.git
cd go-task-scheduler

# Install dependencies
go mod download
```

## Configuration

Create a `.env` file in the root directory (optional):

```env
REDIS_ADDR=localhost:6379  # Redis server address
REDIS_DB=0                 # Redis database number
INTERVAL=200ms             # Worker poll interval
JOB_DELAY=3s              # Default job delay
JOB_MAX_RETRY=3           # Maximum retry attempts
```

If `.env` is not provided, the system uses the defaults shown above.

## Usage

### 1. Start Redis

Make sure Redis is running:

```bash
redis-server
```

### 2. Start the Worker

The worker continuously polls for ready jobs:

```bash
# Using script
./worker.sh

# Or directly
go run cmd/worker/main.go
```

### 3. Schedule Jobs

In a separate terminal, enqueue jobs using the scheduler:

```bash
# Using script (payload, delay in seconds)
./scheduler.sh "Hello World" 5

# Or directly
go run cmd/scheduler/main.go "Hello World" 5

# Default values if arguments not provided
go run cmd/scheduler/main.go
# Uses: payload="Hello from scheduler failed", delay=3s
```

### CLI Arguments

**Scheduler**:
```bash
go run cmd/scheduler/main.go [payload] [delay_in_seconds]
```

- `payload` (optional): Job payload string (default: "Hello from scheduler failed")
- `delay` (optional): Delay in seconds before job execution (default: from `JOB_DELAY` config)

## Examples

### Schedule a successful job

```bash
# Schedule job to run after 3 seconds
./scheduler.sh "Process payment" 3
```

Worker output:
```
✅ Executed job: Process payment
```

### Test retry mechanism

```bash
# Schedule a failing job (payload = "fail")
./scheduler.sh "fail" 2
```

Worker output:
```
⚠️ Job failed: fail (retries left=3)
🔄 Retrying job in 2s (remaining retries=2)
⚠️ Job failed: fail (retries left=2)
🔄 Retrying job in 2s (remaining retries=1)
...
❌ Job exhausted retries: fail
```

### Multiple jobs

```bash
# Terminal 1: Start worker
./worker.sh

# Terminal 2: Schedule multiple jobs
./scheduler.sh "Job 1" 2
./scheduler.sh "Job 2" 5
./scheduler.sh "Job 3" 1
```

## Building

```bash
# Build both binaries
go build -o bin/scheduler cmd/scheduler/main.go
go build -o bin/worker cmd/worker/main.go

# Run built binaries
./bin/worker &
./bin/scheduler "test job" 3
```

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...

# Test specific package
go test ./internal/scheduler
go test ./internal/worker
```

## Project Structure

```
.
├── cmd/
│   ├── scheduler/          # Scheduler entry point
│   │   └── main.go
│   └── worker/             # Worker entry point
│       └── main.go
├── internal/
│   ├── config/             # Configuration loader
│   │   └── config.go
│   ├── job/                # Job struct definition
│   │   └── job.go
│   ├── scheduler/          # Job scheduler logic
│   │   └── scheduler.go
│   └── worker/             # Job worker logic
│       └── worker.go
├── scheduler.sh            # Helper script to run scheduler
├── worker.sh               # Helper script to run worker
├── go.mod
└── README.md
```

## How Retry Works

- Jobs with payload `"fail"` simulate failure
- Failed jobs are retried up to `MaxRetry` times (default: 3)
- Each retry has a 2-second delay
- After exhausting retries, the job is discarded
- Retry count decreases with each attempt

## Dependencies

- [go-redis/v9](https://github.com/redis/go-redis) - Redis client for Go
- [google/uuid](https://github.com/google/uuid) - UUID generation
- [joho/godotenv](https://github.com/joho/godotenv) - Environment variable loading

## License

MIT

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
