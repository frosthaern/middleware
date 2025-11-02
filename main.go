package main

import (
	"fmt"
	middleware "middleware/middleware"
	ratelimiter "middleware/middleware/ratelimiter"
	"net/http"
	"time"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Hello, World ")
}

func main() {
	localhost_addr := ":3000"

	log := middleware.NewLogger("thicka")
	token_bucket := ratelimiter.NewTokenBucket(10, 1, 3 * time.Second)
	fo := middleware.NewMiddlewareFuncBuilder().
		Add(log.Middleware()).
		Add(token_bucket.RateLimiter()).
		Build(handleRoot)

	mux := http.NewServeMux()
	mux.Handle("/", fo)

	if err := http.ListenAndServe(localhost_addr, mux); err != nil {
		fmt.Printf("Server Failed to Start at %s", localhost_addr)
	}
}
