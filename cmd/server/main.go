package main

import (
	"fmt"
	"log"
	"net/http"
	"nft_marketplace/eth/source/handlers"
)

func main() {
    log.SetPrefix("nft_marketplace")
    log.SetFlags(0)

    mux := http.NewServeMux()

    mux.HandleFunc("/hello", handlers.HelloHandler)
    
    fmt.Println("Server starting to listen on port: 8000")
    err := http.ListenAndServe(":8000", mux)

    if err != nil {
        log.Fatal("Server failed to start" , err)
    }
}
