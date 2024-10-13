package users

import (
	"encoding/json"
	"io"
	"net/http"
	"nft_marketplace/eth/source/database"
	"nft_marketplace/eth/source/database/models"
	"nft_marketplace/eth/source/services/users"

	"github.com/gorilla/mux"
)

// type UserDTO struct {
// }

func GetUserByID(w http.ResponseWriter, r *http.Request) {
    muxVars := mux.Vars(r)
    userID := muxVars["id"]

    user, result := users.FindUserById(userID)

    if result.Error != nil {
        http.Error(w, "Record not found", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(user)
}

func CreateNewUser(w http.ResponseWriter, r *http.Request) {
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Error reading body", http.StatusBadRequest)
    }

    var user models.User
    err = json.Unmarshal(body, &user)
    if err != nil {
        http.Error(w, "Error parsing body", http.StatusBadRequest)
    }

    defer r.Body.Close()

    database.Postgres.Create(&user)

    json.NewEncoder(w).Encode(map[string]uint{"id": user.ID})
}
