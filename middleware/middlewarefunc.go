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

type MiddlewareFuncBuilder struct {

	/*
	 * this just adds all the middlewares to a slice
	 * you can add a few mids then do build nad add some more and do build again if you want
	 * but it's usually better to create multiple diff MiddlewareFuncBuilder objects because you may want to add diff middlewares
	 */
	
	middlewares []MiddlewareFunc
}

func NewMiddlewareFuncBuilder() *MiddlewareFuncBuilder {
	return &MiddlewareFuncBuilder{middlewares: []MiddlewareFunc{}}
}

func (mb *MiddlewareFuncBuilder) Add(f MiddlewareFunc) {
	mb.middlewares = append(mb.middlewares, f)
}

func (mb *MiddlewareFuncBuilder) Build(f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	var h http.Handler = http.HandlerFunc(f)
	for i := len(mb.middlewares) - 1; i >= 0;  i-- {
		h = mb.middlewares[i](h)
	}
	return h
}
