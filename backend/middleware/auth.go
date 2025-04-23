package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("Error: Authorization header missing")
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Println("Error: Invalid Authorization header format")
			http.Error(w, "Authorization header must start with Bearer", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		log.Printf("Validating token: %s...", tokenString[:10]) // Log partial token for security

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("Unexpected signing method: %v", token.Header["alg"])
				return nil, http.ErrNotSupported
			}
			return []byte("your-secret-key"), nil
		})

		if err != nil {
			log.Printf("Error parsing token: %v", err)
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			log.Println("Error: Token is invalid or expired")
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Extract username from claims and add to context
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("Error: Invalid token claims")
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}
		username, ok := claims["username"].(string)
		if !ok {
			log.Println("Error: Username not found in token claims")
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Add username to request context
		ctx := context.WithValue(r.Context(), "username", username)
		log.Printf("Authenticated request for user: %s", username)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
