package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/idamidzin/go-task-scheduler/internal/config"
	"github.com/idamidzin/go-task-scheduler/internal/job"
	"github.com/idamidzin/go-task-scheduler/internal/scheduler"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.LoadConfig()

	client := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
		DB:   cfg.RedisDB,
	})

	s := scheduler.NewScheduler(client)

	// Ambil payload dari command line, default "Hello from scheduler failed"
	payload := "Hello from scheduler failed"
	if len(os.Args) > 1 {
		payload = os.Args[1]
	}

	// Ambil delay dari command line, default cfg.JobDelay
	delay := cfg.JobDelay
	if len(os.Args) > 2 {
		if d, err := strconv.Atoi(os.Args[2]); err == nil {
			delay = time.Duration(d) * time.Second
		}
	}

	j := job.Job{
		Payload:  payload,
		Delay:    delay,
		MaxRetry: cfg.JobMaxRetry,
	}

	if err := s.Enqueue(j); err != nil {
		log.Fatalf("❌ failed enqueue: %v", err)
	}
	log.Printf("✅ Job enqueued! (payload=%q, delay=%v, maxRetry=%d)", j.Payload, j.Delay, j.MaxRetry)
}
