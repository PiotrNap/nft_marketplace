package handlers

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(5,1)

type Claims struct {
    Username string `json:"username"`
    jwt.RegisteredClaims
}

func JWTAuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
        // get the string from header
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Missing authorization header", http.StatusUnauthorized)
            return
        }

        // trim the string
        tokenString := strings.TrimPrefix(authHeader , "Bearer ")

        // unpack the Claims 
        claims := &Claims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return token, nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Invalid JWT token", http.StatusUnauthorized)
            return
        }

        next.ServeHTTP(w,r)
    })
}

func RateLimitMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
        
        if !limiter.Allow() {
            http.Error(w, "Max request limit reached", http.StatusTooManyRequests)
        }

        next.ServeHTTP(w,r)
    })
}
