package middleware

import "net/http"

type MiddlewareFunc func(http.Handler) http.Handler

	/*
	 * So in the function make sure to call f.ServeHTTP(r, w)
	 * where r http.Request
	 * and w http.ResponseWriter
	 * and in order to return of type http.Handler use http.HandlerFunc()
	 */

func (m MiddlewareFunc) MiddlewareFu(next http.Handler) http.Handler {
	return m(next)
}

// type MiddlewareFuncBuilder struct {

// 	/*
// 	 * this just adds all the middlewares to a slice
// 	 * you can add a few mids then do build nad add some more and do build again if you want
// 	 * but it's usually better to create multiple diff MiddlewareFuncBuilder objects because you may want to add diff middlewares
// 	 */
	
// 	middlewares []MiddlewareFunc
// }

type Middlewares []MiddlewareFunc

func NewMiddlewareFuncBuilder() Middlewares {
	return make(Middlewares, 0)
}

func (m Middlewares) Add(f MiddlewareFunc) Middlewares {
	return append(m, f)
}

func (m Middlewares) Build(f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	var h http.Handler = http.HandlerFunc(f)
	for i := len(m) - 1; i >= 0;  i-- {
		h = m[i](h)
	}
	return h
}
