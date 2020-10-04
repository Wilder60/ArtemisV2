package middleware

import (
	"net/http"

	"github.com/Wilder60/KeyRing/security"
)

func authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if err := security.Validate(token); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
