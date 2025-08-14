package main

import "net/http"

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; stype-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		w.Header().Set("Referrer-Policy", "orgin-when-cross-orgin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
} //as we want middleware to act on all requests, we need this to exec before request hits our servemux.
//secureHeaders -> servemux -> application handler
