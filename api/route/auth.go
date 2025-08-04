package route

import (
	"net/http"
	"os"
)

var secret = os.Getenv("ADMIN_SECRET")

func withSecret(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const prefix = "Bearer "
		authHeader := r.Header.Get("Authorization")
		if len(authHeader) <= len(prefix) || authHeader[:len(prefix)] != prefix {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		key := authHeader[len(prefix):]
		if key != secret {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		handler(w, r)
	}
}
