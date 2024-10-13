package main

import (
	"fmt"
	"log"
	"net/http"

    "github.com/gorilla/mux"

	// "nft_marketplace/eth/source/handlers"
	// "nft_marketplace/eth/source/handlers/users"
	"nft_marketplace/eth/source/database"
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

    database.Init()

    r := mux.NewRouter()

    // mux.HandleFunc("/hello", handlers.HelloHandler)

    // Users
    r.HandleFunc("/users", users.CreateNewUser).Methods("POST")
    r.HandleFunc("/users/{id}", users.GetUserByID).Methods("GET")
    
    fmt.Println("Server starting to listen on port: 8000")
    err := http.ListenAndServe(":8000", r)

    if err != nil {
        log.Fatal("Server failed to start" , err)
    }
}
