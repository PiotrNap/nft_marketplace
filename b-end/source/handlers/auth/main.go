package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"nft_marketplace/eth/source/services/auth"
	"nft_marketplace/eth/source/services/users"
	"nft_marketplace/eth/source/utils"
	"strconv"
	"strings"
	"time"
)

func CreateNewAccount(w http.ResponseWriter, r *http.Request){
    var dto auth.CreateAccountDTO
    // extract body payload
    data, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Unable to read body object", http.StatusBadRequest)
        return
    }   

    defer r.Body.Close()

    err = json.Unmarshal(data, &dto)
    if err != nil {
        http.Error(w, "Error while parsing body object", http.StatusBadRequest)
        return
    }

    user, err := auth.VerifySignatureAndCreateNewAccount(&dto)
    // generate jwt ...
    if err != nil {
        http.Error(w, err.Error(), http.StatusConflict)
        return
    }

    json.NewEncoder(w).Encode(user)
}

func GenerateChallenge(w http.ResponseWriter, r *http.Request) {
    user, err := utils.ExtractUserFromBody(r)
    if err != nil {
        http.Error(w, "Failed extracting user from body", http.StatusBadRequest)
        return
    }

    _, err = auth.CheckIfUserExists(user.Username)
    if user.Verified {
        http.Error(w, "User with a given username already exists", http.StatusBadRequest)
        return
    }
    requestID := r.Header.Get("X-Request-ID")

    if user.Challenge != "" {
        parts := strings.Split(user.Challenge, "_")
        oldChallenge := parts[0]
        fmt.Println("Old challenge :" + oldChallenge)

        ttl, err := strconv.Atoi(oldChallenge)
        if err != nil {
            http.Error(w, "Problem occured on the server", http.StatusBadRequest)
            msg := fmt.Errorf("Error during user.Challenge convertion: %s", err)
            fmt.Println(msg)
            return
        }
        now :=  time.Now().Unix()

        if int(now) < ttl {
            hashedChallenge := utils.HashValues([]byte(oldChallenge))
            json.NewEncoder(w).Encode(map[string]string{"challenge": hashedChallenge})
        }
    } else {
        challenge, err := utils.GenerateChallengeString()
        hashedChallenge := utils.HashValues([]byte(challenge))
        fmt.Println("hashedChallenge ", hashedChallenge)

        if err != nil {
            print(err)
            http.Error(w, "Something went wrong generating challenge string", http.StatusConflict)
            return
        }


        user.Challenge = challenge
        user.Salt = requestID
        users.AddUser(&user)

        json.NewEncoder(w).Encode(map[string]string{"challenge": hashedChallenge})
    }
}

