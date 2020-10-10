package middleware

import (
	"net/http"

	"github.com/Wilder60/KeyRing/internal/security"
)

// Authorize will check the authorization for a given request, this will check if they just have a valid token
// not if they are an adminstrator, that will be handed by the AuthorizeAdmin function
func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		err := security.Validate(token)
		if err == security.ErrInvalidToken {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		} else if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
