package middleware

import (
	"net/http"
	"strings"
)

func Sanitize(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") || (!strings.Contains(r.URL.Path, ".")) {
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}
