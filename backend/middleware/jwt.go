package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/harip/GoTasker/config"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("Error: Authorization header is missing")
			http.Error(w, `{"error": "Authorization header is required"}`, http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Println("Error: Invalid token format, missing Bearer prefix")
			http.Error(w, `{"error": "Invalid token format"}`, http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("Error: Unexpected signing method: %v", token.Header["alg"])
				return nil, http.ErrAbortHandler
			}
			return []byte(config.AppConfig.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			log.Printf("Error: Invalid or expired token: %v", err)
			http.Error(w, `{"error": "Invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			var userID float64
			switch id := claims["user_id"].(type) {
			case float64:
				userID = id
			case int:
				userID = float64(id)
			default:
				log.Println("Error: user_id in token claims is not a number")
				http.Error(w, `{"error": "Invalid token claims"}`, http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "user_id", userID)
			r = r.WithContext(ctx)
		} else {
			log.Println("Error: Invalid token claims format")
			http.Error(w, `{"error": "Invalid token claims"}`, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
