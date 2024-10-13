package main

import (
	"fmt"
	"log"
	"net/http"
	"nft_marketplace/eth/source/handlers"
	"nft_marketplace/eth/source/db"
)

// store nfts info hash with:
// - owners' address
// - list of current bid
// - `itemInfoHash` (value used inside smart contract)
// - duration of the auction

func main() {
    log.SetPrefix("nft_marketplace")
    log.SetFlags(0)

    db.Init()

    mux := http.NewServeMux()

    mux.HandleFunc("/hello", handlers.HelloHandler)
    
    fmt.Println("Server starting to listen on port: 8000")
    err := http.ListenAndServe(":8000", mux)

    if err != nil {
        log.Fatal("Server failed to start" , err)
    }
}
