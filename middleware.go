package main

import (
	"log"
	"net/http"
)

// LogRequest is a middleware that just logs a request
func LogRequest(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		pattern := `%s - "%s %s %s"`
		log.Printf(pattern, r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// SecureHeaders is a middleware that sets the following security related HTTP Headers:
// X-Content-Type-Options: nosniff
// X-Frame-Options: deny
// X-XSS-Protection: 1; mode=block
func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header()["X-XSS-Protection"] = []string{"1; mode=block"}

		next.ServeHTTP(w, r)
	})
}
