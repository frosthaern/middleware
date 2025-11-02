package middleware

import (
	"log"
	"net/http"
	"os"
	"time"
)

type Logger struct {
	prefix string
	logger *log.Logger
}

func NewLogger(prefix string) *Logger {
	return &Logger{
		prefix: prefix,
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (l *Logger) Middleware() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			l.logger.Printf("%s %s %s %v",
				l.prefix,
				r.Method,
				r.URL.Path,
				time.Since(start),
			)
		})
	}
}
