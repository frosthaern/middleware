package ratelimiter

import (
	middleware "middleware/middleware"
	"net/http"
	"sync"
	"time"
)

type TokenBucket struct {
	capacity     int           // this is max no of tokens at any point of time
	tokens       int           // no of currently available tokens
	refillRate   int           // how many tokens to fill whenever a period happens
	refillPeriod time.Duration // time interval between each refill
	lastRefill   time.Time
	mu           sync.Mutex
}

func NewTokenBucket(capacity int, refillRate int, refillPeriod time.Duration) *TokenBucket {
	return &TokenBucket{
		capacity: capacity,
		refillRate: refillRate,
		refillPeriod: refillPeriod,
		tokens: capacity,
		lastRefill: time.Now(),
	}
	
}

func (t *TokenBucket) RateLimiter() middleware.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !t.Allow() {
				http.Error(w, "too many requests per second", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// refill()
/*
 * i hope that when it is initialized it has all the tokens
 * and the last refill time is set
 * so that i can always call the refill and it's valid
 */

func (t *TokenBucket) refill() {
	prev := t.lastRefill
	now := time.Now()
	periods := int(now.Sub(prev) / t.refillPeriod)

	if periods > 0 {
		t.tokens += t.refillRate * periods
		t.tokens = min(t.capacity, t.tokens)

		// because if you are in between two time periods, you only fill up till last period
		// and that's why your lastRefill will also be the last period

		t.lastRefill = t.lastRefill.Add(time.Duration(periods) * t.refillPeriod)
	}
}

func (t *TokenBucket) Allow() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.refill()

	if t.tokens > 0 {
		t.tokens--
		return true
	}
	return false
}

func (t *TokenBucket) GetRemainingTokens() int {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.refill()

	return t.tokens
}
