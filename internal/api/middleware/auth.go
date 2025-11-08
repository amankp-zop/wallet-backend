package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDContextKey = contextKey("userID")

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			authHeader := r.Header.Get("Authorization")
			if authHeader == ""{
				http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
				
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts)<2 || parts[0]!="Bearer"{
				http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
				
				return
			}

			tokenString := parts[1]

			token,err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok{
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				
				return []byte(jwtSecret), nil
			})
			if err!=nil{
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid{
				userIDFloat, ok:= claims["sub"].(float64)
				if !ok{
					http.Error(w, "Invalid token claims", http.StatusUnauthorized)
					
					return
				}

				userID := int64(userIDFloat)

				ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
				next.ServeHTTP(w, r.WithContext(ctx))
			}else{
				http.Error(w, "Invalid token", http.StatusUnauthorized)
			}

		})
	}
}