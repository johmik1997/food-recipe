package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `{"error":"internal server error"}`)
				debug.PrintStack() // Log stack trace
			}
		}()
		next.ServeHTTP(w, r)
	})
}