package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"swiftShare/backend/controllers"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserKey contextKey = "user"

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie("jwt")
		if err != nil {
			http.Error(w, "Could not get token information.", http.StatusUnauthorized)
			return
		}
		token, err := jwt.Parse(tokenCookie.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("SECRET_KEY")), nil
		})
		if err != nil {
			http.Error(w, "Error parsing JWT.", http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			http.Error(w, "Invalid JWT.", http.StatusUnauthorized)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || float64(time.Now().Unix()) > claims["exp"].(float64) {
			http.Error(w, "Expired JWT.", http.StatusUnauthorized)
			return
		}
		userID, ok := claims["sub"]
		if !ok {
			http.Error(w, "Invalid user ID in JWT.", http.StatusUnauthorized)
			return
		}
		user, err := controllers.FindByID(fmt.Sprint(userID))
		if err != nil {
			http.Error(w, "Error retrieving user information.", http.StatusUnauthorized)
			log.Println(err)
			return
		}
		ctx := context.WithValue(r.Context(), UserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
