package main

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type TestResult struct {
	successCount int32
	errorCount   int32
	totalTime    int64
	mu           sync.Mutex
	requestTypes map[string]*RequestStats
}

type RequestStats struct {
	Count     int32
	Success   int32
	Errors    int32
	TotalTime int64
	MinTime   int64
	MaxTime   int64
}

func (r *TestResult) recordRequest(reqType string, success bool, duration int64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	stats, exists := r.requestTypes[reqType]
	if !exists {
		stats = &RequestStats{MinTime: 1<<63 - 1}
		r.requestTypes[reqType] = stats
	}

	stats.Count++
	if success {
		stats.Success++
		stats.TotalTime += duration
		if duration < stats.MinTime {
			stats.MinTime = duration
		}
		if duration > stats.MaxTime {
			stats.MaxTime = duration
		}
	} else {
		stats.Errors++
	}

	if success {
		atomic.AddInt32(&r.successCount, 1)
	} else {
		atomic.AddInt32(&r.errorCount, 1)
	}
	atomic.AddInt64(&r.totalTime, duration)
}

func (r *TestResult) printStats(duration time.Duration, total int) {
	fmt.Printf("\nSummary:\n")
	fmt.Printf("Total Requests:    %d\n", total)
	fmt.Printf("Successful:        %d (%.1f%%)\n", r.successCount, float64(r.successCount)/float64(total)*100)
	fmt.Printf("Errors:            %d (%.1f%%)\n", r.errorCount, float64(r.errorCount)/float64(total)*100)
	fmt.Printf("Total Duration:    %.2f seconds\n", duration.Seconds())
	fmt.Printf("Requests/sec:      %.2f\n", float64(total)/duration.Seconds())
	fmt.Printf("Avg. Latency:      %.2f ms\n", float64(atomic.LoadInt64(&r.totalTime))/float64(total))

	fmt.Printf("\nRequest Types Breakdown:\n")
	fmt.Printf("%-25s %8s %8s %8s %10s %10s %10s\n",
		"Type", "Total", "Success", "Errors", "Success %", "Min(ms)", "Max(ms)")

	r.mu.Lock()
	for reqType, stats := range r.requestTypes {
		successPct := float64(stats.Success) / float64(stats.Count) * 100

		fmt.Printf("%-25s %8d %8d %8d %9.1f%% %9d %9d\n",
			reqType, stats.Count, stats.Success, stats.Errors, successPct, stats.MinTime, stats.MaxTime)
	}
	r.mu.Unlock()
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	}))
	slog.SetDefault(logger)

	cfg := LoadConfig()

	initializeTestData(cfg)

	result := &TestResult{
		requestTypes: make(map[string]*RequestStats),
	}
	sem := make(chan struct{}, cfg.Concurrency)
	var wg sync.WaitGroup

	fmt.Printf("Starting load test...\n")
	fmt.Printf("Target: %s\n", cfg.BaseURL)
	fmt.Printf("Requests: %d, Concurrency: %d\n\n", cfg.TotalRequests, cfg.Concurrency)

	start := time.Now()

	for i := 0; i < cfg.TotalRequests; i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func(id int) {
			defer wg.Done()
			defer func() { <-sem }()

			reqType, request, expectedStatus := generateRequest(id, cfg.BaseURL)
			requestStart := time.Now()

			err := sendRequest(cfg, request, expectedStatus)
			duration := time.Since(requestStart).Milliseconds()

			result.recordRequest(reqType, err == nil, duration)

			if err != nil {
				if !isExpectedError(err) {
					slog.Warn("Request failed", "type", reqType, "error", err)
				}
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)
	result.printStats(duration, cfg.TotalRequests)
}
