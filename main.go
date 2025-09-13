package main

import (
	"fmt"
	"net/http"
	middleware "middleware/middleware"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Hello, World ")
}

func main() {
	localhost_addr := ":3000"

	log := middleware.NewLogger("thicka")
	fo := middleware.NewMiddlewareFuncBuilder().Add(log.Middleware()).Build(handleRoot)

	mux := http.NewServeMux()
	mux.Handle("/", fo)
	
	if err := http.ListenAndServe(localhost_addr, mux); err != nil {
		fmt.Printf("Server Failed to Start at %s", localhost_addr)
	}
}
