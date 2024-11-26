package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"nft_marketplace/eth/source/database"
	"nft_marketplace/eth/source/handlers"
	"nft_marketplace/eth/source/handlers/auth"
	"nft_marketplace/eth/source/handlers/users"
)

// store nfts info hash with:
// - owners' address
// - list of current bid
// - `itemInfoHash` (value used inside smart contract)
// - duration of the auction

func main() {
    log.SetPrefix("nft_marketplace")
    log.SetFlags(0)

    if jwtSecret :=  os.Getenv("JWT_SECRET"); jwtSecret == "" {
       log.Fatal("Missing env variable JWT_SECRET") 
    }

    database.Init()

    r := mux.NewRouter()

    // Users
    r.HandleFunc("/users", users.CreateNewUser).Methods("POST")
    r.Handle("/users/{id}",
        handlers.JWTAuthMiddleware(http.HandlerFunc(users.GetUserByID))).Methods("GET")
    r.HandleFunc("/users/username", users.CheckIfUserExists).Methods("POST")

    // Authentication
    r.HandleFunc("/auth/challenge", auth.GenerateChallenge).Methods("POST") 
    // r.Handle("/auth/challenge", handlers.RequestIDMiddleware(http.HandlerFunc(auth.GenerateChallenge))).Methods("POST")

    r.HandleFunc("/auth/sign-up", auth.CreateNewAccount).Methods("POST") 
    // r.HandleFunc("/auth/challenge/verify", auth.VerifyChallenge).Methods("POST")

    c := cors.New(cors.Options{
        AllowedOrigins:     []string{"http://localhost:3000"},
        AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE"},
        AllowedHeaders:     []string{"Authorization", "Content-Type"},
        AllowCredentials:   true,
    })

    handler := handlers.RequestIDMiddleware(handlers.RateLimitMiddleware(c.Handler(r)))

    
    fmt.Println("Server starting to listen on port: 8000")
    err := http.ListenAndServe(":8000", handler)

    if err != nil {
        log.Fatal("Server failed to start" , err)
    }
}
