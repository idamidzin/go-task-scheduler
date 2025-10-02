package scheduler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/idamidzin/go-task-scheduler/internal/job"
	"github.com/redis/go-redis/v9"
)

type Scheduler struct {
	client *redis.Client
}

func NewScheduler(client *redis.Client) *Scheduler {
	return &Scheduler{client: client}
}

func (s *Scheduler) Enqueue(j job.Job) error {
	if j.ID == "" {
		j.ID = uuid.New().String()
	}

	score := float64(time.Now().Add(j.Delay).UnixMilli())

	data, err := json.Marshal(j)
	if err != nil {
		return err
	}

	return s.client.ZAdd(context.Background(), "jobs", redis.Z{
		Score:  score,
		Member: data,
	}).Err()
}
