package main

import (
	"net/http"
	"strings"
)

func ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		const BEARER string = "Bearer "
		bearerToken := r.Header.Get("Authorization")
		if len(bearerToken) <= 0 || !strings.HasPrefix(bearerToken, BEARER) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(http.StatusText(http.StatusNotFound)))
			return
		}
		token := strings.SplitAfter(bearerToken, BEARER)[0]
		ValidateJWT(token)
		next.ServeHTTP(w, r)
	})
}
