package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	RedisAddr   string
	RedisDB     int
	Interval    time.Duration
	JobDelay    time.Duration
	JobMaxRetry int
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	db, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

	interval, err := time.ParseDuration(getEnv("INTERVAL", "200ms"))
	if err != nil {
		interval = 200 * time.Millisecond
	}

	delay, err := time.ParseDuration(getEnv("JOB_DELAY", "3s"))
	if err != nil {
		delay = 3 * time.Second
	}

	maxRetry, err := strconv.Atoi(getEnv("JOB_MAX_RETRY", "3"))
	if err != nil {
		maxRetry = 3
	}

	return &Config{
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
		RedisDB:     db,
		Interval:    interval,
		JobDelay:    delay,
		JobMaxRetry: maxRetry,
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
