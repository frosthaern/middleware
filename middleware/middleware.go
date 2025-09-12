package middleware

import "net/http"

// type Technique int

// const (
// 	SlidingWindowLog Technique = iota
// 	TokenBucket
// 	SlidingWindowCounter
// 	FixedWindowCounter
// 	LeakyBucket
// )

type MiddleWare struct {
	route map[string]http.Handler
}

func (m *MiddleWare) Handle(path string, f http.Handler) {
	m.route[path] = f
}

func NewMiddleWare() *MiddleWare {
	return &MiddleWare{route: make(map[string]http.Handler)}
}

func (m *MiddleWare) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// this is where the main logic is going to be
	// you need to call the f.ServeHTTP(w http.ResponseWriter, r *http.Request),
	// so that function can run, but before that you can put some other logic and stuff
	url_path := r.URL.Path
	f, ok := m.route[url_path]
	if ok {
		f.ServeHTTP(w, r)
	}
}

// now you have to implement ServeHTTP function and while doing that you can do a lot more middleware stuff
/*
 * like you will have middleware functions and values in the MiddleWare struct
 * and in the ServeHTTP function we can either call those functions or pass those values and stuff and then we can use those in the ServeHTTP() function anyways because it uses MiddleWare strcut values
 * tomorrow morning wake up and
 * and you need middleware.AddMiddlewareFunc() with some sig
 */

// i need a function that applies middleware
// and then a function that takes the actual function that takes func ( w http.ResponseWriter, r *http.Request )

type MiddleWareApplier interface {
	RouteHandler(fn func(w http.ResponseWriter, r *http.Request), hs ...func(http.Handler) http.Handler) http.Handler
}
