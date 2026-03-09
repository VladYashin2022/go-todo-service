package httpServer

import (
	"log"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		path := r.URL.Path
		if r.URL.RawQuery != "" {
			path += "?" + r.URL.RawQuery
		}

		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start)

		log.Printf("%s %s %s", r.Method, path, duration)
	})
}