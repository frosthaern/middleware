package ratelimiter

import (
	"middleware/middleware"
	"net/http"
	"sync"
	"time"
)

type LeakyBucket struct {
	capacity    int
	leakRate    time.Duration
	queue       chan struct{}
	mu          sync.Mutex
	lastLeak    time.Time
	stopLeaking chan struct{}
}

// NewLeakyBucket creates a new LeakyBucket rate limiter
// capacity: maximum number of requests that can be in the bucket
// leakRate: how often a request is allowed to leak through (e.g., 100ms means 10 requests per second)
func NewLeakyBucket(capacity int, leakRate time.Duration) *LeakyBucket {
	lb := &LeakyBucket{
		capacity:    capacity,
		leakRate:    leakRate,
		queue:       make(chan struct{}, capacity),
		lastLeak:    time.Now(),
		stopLeaking: make(chan struct{}),
	}
	go lb.leak()
	return lb
}

func (lb *LeakyBucket) leak() {
	ticker := time.NewTicker(lb.leakRate)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			lb.mu.Lock()
			select {
			case <-lb.queue: // Remove an item from the queue if available
			default: // Do nothing if queue is empty
			}
			lb.lastLeak = time.Now()
			lb.mu.Unlock()
		case <-lb.stopLeaking:
			return
		}
	}
}

func (lb *LeakyBucket) Allow() bool {
	select {
	case lb.queue <- struct{}{}:
		return true
	default:
		return false
	}
}

func (lb *LeakyBucket) Stop() {
	close(lb.stopLeaking)
}

func (lb *LeakyBucket) RateLimiter() middleware.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !lb.Allow() {
				http.Error(w, "too many requests", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
