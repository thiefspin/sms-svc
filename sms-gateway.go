package main

import (
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

const LIMIT_PER_MINUTE = 100

var limiter = rate.NewLimiter(rate.Every(1*time.Minute/LIMIT_PER_MINUTE), LIMIT_PER_MINUTE)

func ratelimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Sent"))
	})
	http.ListenAndServe(":4000", ratelimit(mux))
}
