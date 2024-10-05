package handlers

import (
	"fmt"
	"net/http"
	"nft_marketplace/eth/source/services"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
    key := "VERYSecretKey"

    fmt.Println("Hello from the server! Getting to work now...")

    msg, err := services.HelloService(key)
    if err != nil {
        http.Error(w, "We're sorry, something went wrong on our side.", http.StatusInternalServerError)
    }
    fmt.Fprintln(w, msg)
}
