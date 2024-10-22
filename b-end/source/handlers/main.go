package handlers

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    Username string `json:"username"`
    jwt.RegisteredClaims
}

func GenerateJWT(username string) (string, error) {
    var jwtSecretString = os.Getenv("JWT_SECRET")
    jwtSecretBytes := []byte(jwtSecretString)

    expirationTime := time.Now().Add(time.Hour * 24)

    claims := &Claims{
        Username: username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    tokenString, err := token.SignedString(jwtSecretBytes)
    if err != nil {
        return "", err
    }

    return tokenString, nil
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
