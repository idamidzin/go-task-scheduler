package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/idamidzin/go-task-scheduler/internal/job"
	"github.com/redis/go-redis/v9"
)

type Worker struct {
	client   *redis.Client
	interval time.Duration
}

func NewWorker(client *redis.Client, interval time.Duration) *Worker {
	return &Worker{client: client, interval: interval}
}

func (w *Worker) Start() {
	log.Println("🚀 Worker started...")
	ctx := context.Background()

	for {
		now := float64(time.Now().UnixMilli())

		res, err := w.client.ZRangeByScore(ctx, "jobs", &redis.ZRangeBy{
			Min:    "-inf",
			Max:    fmt.Sprintf("%f", now),
			Offset: 0,
			Count:  1,
		}).Result()

		if err != nil {
			log.Printf("❌ error reading jobs: %v", err)
			time.Sleep(w.interval)
			continue
		}

		if len(res) == 0 {
			log.Printf("⏳ Worker idle, no job ready (check every %v)", w.interval)
			time.Sleep(w.interval)
			continue
		}

		rawJob := res[0]

		// Decode job
		var j job.Job
		if err := json.Unmarshal([]byte(rawJob), &j); err != nil {
			log.Printf("❌ failed to decode job: %v", err)
			log.Printf("raw job data: %s", rawJob)
			// Hapus job yang corrupt supaya tidak loop error terus
			_, _ = w.client.ZRem(ctx, "jobs", rawJob).Result()
			continue
		}

		// Remove from queue (pop)
		_, _ = w.client.ZRem(ctx, "jobs", rawJob).Result()

		// Simulasi eksekusi job → gagal kalau payload "fail"
		if j.Payload == "fail" {
			log.Printf("⚠️ Job failed: %s (retries left=%d)", j.Payload, j.MaxRetry)
			if j.MaxRetry > 0 {
				// retry dengan delay 2 detik
				j.MaxRetry--
				j.Delay = 2 * time.Second

				data, err := json.Marshal(j)
				if err != nil {
					log.Printf("❌ failed to marshal job for retry: %v", err)
					continue
				}

				score := float64(time.Now().Add(j.Delay).UnixMilli())
				if _, err := w.client.ZAdd(ctx, "jobs", redis.Z{Score: score, Member: data}).Result(); err != nil {
					log.Printf("❌ failed to push job to Redis for retry: %v", err)
					continue
				}

				log.Printf("🔄 Retrying job in %v (remaining retries=%d)", j.Delay, j.MaxRetry)
			} else {
				log.Printf("❌ Job exhausted retries: %s", j.Payload)
			}
			continue
		}

		// Kalau sukses
		log.Printf("✅ Executed job: %s", j.Payload)
	}
}
