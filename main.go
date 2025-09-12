package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
	mid "rate-limiter/middleware"
)

var (
	route_time_map = make(map[string]time.Time)
	mu             sync.RWMutex
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Hello, World")
}

func main() {
	localhost_addr := ":3000"

	
	middleware := mid.NewMiddleWare()
	middleware.Handle("/", http.HandlerFunc(handleRoot))

	
	if err := http.ListenAndServe(localhost_addr, middleware); err != nil {
		fmt.Printf("Server Failed to Start at %s", localhost_addr)
	}
}

func middleware1(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("printing first middleware shit")
		f.ServeHTTP(w, r)
	})
}


func middleware2(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("printing second middleware shit")
		f.ServeHTTP(w, r)
	})
}


func middleware3(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("printing third middleware shit")
		f.ServeHTTP(w, r)
	})
}
