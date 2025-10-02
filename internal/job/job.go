package job

import "time"

type Job struct {
	ID       string
	Payload  string
	Delay    time.Duration
	Retries  int
	MaxRetry int
}
