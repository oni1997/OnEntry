package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/oni1997/onentry/services/api-go/database"
	"github.com/oni1997/onentry/services/api-go/models"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthMiddleware(db *database.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				models.JSONError(w, http.StatusUnauthorized, "Missing authorization header")
				return
			}

			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				models.JSONError(w, http.StatusUnauthorized, "Invalid authorization header format")
				return
			}

			tokenHash := tokenParts[1]
			session, err := db.GetSessionByTokenHash(r.Context(), tokenHash)
			if err != nil {
				models.JSONError(w, http.StatusUnauthorized, "Invalid token")
				return
			}

			if session.ExpiresAt.Before(models.Now()) {
				_ = db.DeleteSession(r.Context(), session.ID)
				models.JSONError(w, http.StatusUnauthorized, "Token expired")
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserFromContext(ctx context.Context) *models.Session {
	if session, ok := ctx.Value(UserContextKey).(*models.Session); ok {
		return session
	}
	return nil
}
