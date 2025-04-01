package main

import (
	"net/http"
	"strings"
)

func VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const BEARER string = "Bearer "
		bearerToken := r.Header.Get("Authorization")
		if len(bearerToken) <= 0 || !strings.HasPrefix(bearerToken, BEARER) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(http.StatusText(http.StatusNotFound)))
			return
		}
		token := strings.SplitAfter(bearerToken, BEARER)[1]
		if _, err := VerifyJWT(token); err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(http.StatusText(http.StatusNotFound)))
			return
		}
		next.ServeHTTP(w, r)
	})
}
