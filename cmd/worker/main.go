package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/idamidzin/go-task-scheduler/internal/config"
	"github.com/idamidzin/go-task-scheduler/internal/worker"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.LoadConfig()

	client := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
		DB:   cfg.RedisDB,
	})

	w := worker.NewWorker(client, cfg.Interval)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go w.Start()

	<-stop
}
