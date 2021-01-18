package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
)

func JSONMiddleWare() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Set("Content-Type", "application/json;charset=utf-8")

			next.ServeHTTP(rw, r)
		})
	}
}
