package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func TimerMiddleWare() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(rw, r)

			passed := time.Now().Sub(start)

			log.Printf("Request executed in %d ms\n", passed.Milliseconds())
		})
	}
}
