package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
	middleware "rate-limiter/middleware"
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

	
	mid := middleware.NewMiddlewareFuncBuilder()
	mid.Add(middleware1)
	mid.Add(middleware2)
	mid.Add(middleware3)
	overall_func := mid.Build(handleRoot)

	// this is how i want the feature to be
	root := middleware.NewMiddlewareRouteBuilder()
	root.Set("/", overall_func)

	
	if err := http.ListenAndServe(localhost_addr, root); err != nil {
		fmt.Printf("Server Failed to Start at %s", localhost_addr)
	}
}

func middleware1(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("printing first middleware shit\n")
		f.ServeHTTP(w, r)
	})
}


func middleware2(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("printing second middleware shit\n")
		f.ServeHTTP(w, r)
	})
}


func middleware3(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("printing third middleware shit\n")
		f.ServeHTTP(w, r)
	})
}
