package middleware

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type metricsStore struct {
	mu sync.Mutex

	// key: method|path|status
	reqCount map[string]uint64

	// key: method|path
	durCount map[string]uint64
	durSum   map[string]float64
}

var globalMetrics = &metricsStore{
	reqCount: make(map[string]uint64),
	durCount: make(map[string]uint64),
	durSum:   make(map[string]float64),
}

type metricsStatusWriter struct {
	http.ResponseWriter
	status int
}

func (w *metricsStatusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// Metrics collects basic HTTP counters and durations in-memory.
func Metrics() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sw := &metricsStatusWriter{ResponseWriter: w, status: http.StatusOK}
			start := time.Now()

			next.ServeHTTP(sw, r)

			method := r.Method
			path := r.URL.Path
			status := sw.status

			dur := time.Since(start).Seconds()

			globalMetrics.mu.Lock()
			globalMetrics.reqCount[fmt.Sprintf("%s|%s|%d", method, path, status)]++
			key := fmt.Sprintf("%s|%s", method, path)
			globalMetrics.durCount[key]++
			globalMetrics.durSum[key] += dur
			globalMetrics.mu.Unlock()
		})
	}
}

// MetricsHandler exposes metrics in Prometheus text exposition format.
func MetricsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")

		globalMetrics.mu.Lock()
		// copy for deterministic output
		reqCopy := make(map[string]uint64, len(globalMetrics.reqCount))
		for k, v := range globalMetrics.reqCount {
			reqCopy[k] = v
		}
		durCountCopy := make(map[string]uint64, len(globalMetrics.durCount))
		for k, v := range globalMetrics.durCount {
			durCountCopy[k] = v
		}
		durSumCopy := make(map[string]float64, len(globalMetrics.durSum))
		for k, v := range globalMetrics.durSum {
			durSumCopy[k] = v
		}
		globalMetrics.mu.Unlock()

		fmt.Fprintln(w, "# HELP http_requests_total Total number of HTTP requests.")
		fmt.Fprintln(w, "# TYPE http_requests_total counter")

		reqKeys := make([]string, 0, len(reqCopy))
		for k := range reqCopy {
			reqKeys = append(reqKeys, k)
		}
		sort.Strings(reqKeys)

		for _, k := range reqKeys {
			parts := strings.Split(k, "|")
			method, path, statusStr := parts[0], parts[1], parts[2]
			status, _ := strconv.Atoi(statusStr)
			fmt.Fprintf(w, "http_requests_total{method=%q,path=%q,status=%q} %d\n", method, path, strconv.Itoa(status), reqCopy[k])
		}

		fmt.Fprintln(w, "# HELP http_request_duration_seconds_sum Total time spent serving requests (seconds).")
		fmt.Fprintln(w, "# TYPE http_request_duration_seconds_sum counter")
		fmt.Fprintln(w, "# HELP http_request_duration_seconds_count Total number of observed request durations.")
		fmt.Fprintln(w, "# TYPE http_request_duration_seconds_count counter")

		durKeys := make([]string, 0, len(durCountCopy))
		for k := range durCountCopy {
			durKeys = append(durKeys, k)
		}
		sort.Strings(durKeys)

		for _, k := range durKeys {
			parts := strings.Split(k, "|")
			method, path := parts[0], parts[1]
			fmt.Fprintf(w, "http_request_duration_seconds_sum{method=%q,path=%q} %.6f\n", method, path, durSumCopy[k])
			fmt.Fprintf(w, "http_request_duration_seconds_count{method=%q,path=%q} %d\n", method, path, durCountCopy[k])
		}
	})
}
