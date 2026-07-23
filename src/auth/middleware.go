package auth

import (
	"log/slog"
	"net/http"
	"strings"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("auth middleware start method=%s path=%s", r.Method, r.URL.Path)

		// Get header from the request
		authHeader := r.Header.Get("Authorization")

		// If header is empty, unauthorized
		if authHeader == "" {
			slog.Error("auth middleware missing authorization header method=%s path=%s", r.Method, r.URL.Path)
			http.Error(w, "Missing auth header", http.StatusUnauthorized)
			return
		}

		// Split the header value into 2 parts. Ex: header value -> `Bearer 2xcsfseslkjgsjdf...`
		// Result (stored in parts) -> ["Bearer", "2xcsfseslkjgsjdf..."]
		parts := strings.SplitN(authHeader, " ", 2)

		// log.Printf("Splitted authheader: %v", parts)

		// If header value has more than 2 parts or if the first part is not Bearer, invalid value -> unauthorized
		if len(parts) != 2 || parts[0] != "Bearer" {
			slog.Error("auth middleware invalid authorization header method=%s path=%s", r.Method, r.URL.Path)
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		// Get the token value from the header value
		tokenString := parts[1]

		// Parse the token to get the claims
		claims, err := ParseToken(tokenString)
		if err != nil {
			slog.Error("auth middleware token parse failed method=%s path=%s err=%v", r.Method, r.URL.Path, err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if claims == nil {
			slog.Error("auth middleware empty claims method=%s path=%s", r.Method, r.URL.Path)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		slog.Info("auth middleware authorized method=%s path=%s user_id=%s", r.Method, r.URL.Path, claims.UserID)

	})
}
