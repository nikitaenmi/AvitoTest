package main

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	BaseURL       string
	TotalRequests int
	Concurrency   int
	Timeout       time.Duration
}

func LoadConfig() *Config {
	_ = godotenv.Load(".env")

	totalRequests, _ := strconv.Atoi(getEnv("TOTAL_REQUESTS", "500"))
	concurrency, _ := strconv.Atoi(getEnv("CONCURRENCY", "20"))
	timeout, _ := time.ParseDuration(getEnv("TIMEOUT", "10s"))

	return &Config{
		BaseURL:       getEnv("BASE_URL", "http://localhost:8080"),
		TotalRequests: totalRequests,
		Concurrency:   concurrency,
		Timeout:       timeout,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
