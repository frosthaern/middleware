package main

import (
	"fmt"
	"time"
	"sync"
	"net/http"
)

var (
	
	route_time_map = make(map[string]time.Time)
	mu sync.RWMutex
)


func handleRoot(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	var req_url_path = r.URL.Path
	value, ok := route_time_map[req_url_path]
	var this_time = time.Now()
	fmt.Printf("request obj: %s\n", req_url_path);
	if ok {
		var lapsed_time = this_time.Second() - value.Second()
		if lapsed_time < 5 {
			fmt.Printf("request was rejected at second %d ", time.Now().Second())
			return
		} else {
			route_time_map[req_url_path] = this_time
			fmt.Printf("this request waxs accepted at second %d", this_time.Second())
			return
		}
	} else {
		route_time_map[req_url_path] = this_time
		fmt.Printf("this request waxs accepted at second %d", this_time.Second())
		return		
	}
}

func main() {
	localhost_addr := ":3000"
	http.HandleFunc("/", handleRoot);
	if err := http.ListenAndServe(localhost_addr, nil); err != nil {
		fmt.Printf("Server Failed to Start at %s", localhost_addr)
	}
}
