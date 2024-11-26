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

type UserDTO struct {
    Username string `gorm:"json:stirng"`
}

func CheckIfUserExists(w http.ResponseWriter, r *http.Request) {
   var userDTO UserDTO

   body, err := io.ReadAll(r.Body)
   if err != nil {
       http.Error(w, "Error reading body", http.StatusBadRequest)
       return
   }
   defer r.Body.Close()

   if err := json.Unmarshal(body, &userDTO); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
   }

   _, result := users.FindUserByUsername(userDTO.Username)

   print(result)

   if result.Error != nil {
       json.NewEncoder(w).Encode(map[string]bool{"exists": false})
   } else {
       json.NewEncoder(w).Encode(map[string]bool{"exists": true})
   }
}

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
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    defer r.Body.Close()

    database.Postgres.Create(&user)

    json.NewEncoder(w).Encode(map[string]uint{"id": user.ID})
}
