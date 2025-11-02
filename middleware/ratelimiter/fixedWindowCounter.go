package ratelimiter

import (
	"middleware/middleware"
	"net/http"
	"sync"
	"time"
)

type FixedWindowCounter struct {
	slotDuration time.Duration // time duration per slot
	capacity int // now of requests you can take per slot
	processed int // no of requests alread processed in this timeslot
	lastProcessed time.Time
	mu sync.Mutex
}

func (fwc *FixedWindowCounter) NewFixedWindowCounter(slotDuration time.Duration, capacity int) *FixedWindowCounter {
	return &FixedWindowCounter{
		slotDuration: slotDuration,
		capacity: capacity,
		processed: 0,
		lastProcessed: time.Now(),
	}
}

func (fwc *FixedWindowCounter) RateLimiter() middleware.MiddlewareFunc {
	return func(f http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !fwc.Allow() {
				http.Error(w, "too many requests per second", http.StatusTooManyRequests)
			}
			f.ServeHTTP(w, r)
		})
	}
}

func (fwc *FixedWindowCounter) Allow() bool {
	fwc.mu.Lock()
	defer fwc.mu.Unlock()
	now := time.Now()
	period_diff := int(now.Sub(fwc.lastProcessed) / fwc.slotDuration)

	if period_diff > 0 {
		fwc.lastProcessed = fwc.lastProcessed.Add(time.Duration(period_diff) * fwc.slotDuration)
		return true
	}
	
	if fwc.processed < fwc.capacity {
		fwc.processed += 1
		return true
	}
	return false
}
