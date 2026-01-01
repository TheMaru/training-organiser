package auth

import (
	"context"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type contextKey string

const UserIDKey contextKey = "userID"

func MiddlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := GetBearerToken(r.Header)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
