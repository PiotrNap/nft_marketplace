package users

import (
	"nft_marketplace/eth/source/database"
	"nft_marketplace/eth/source/database/models"

	"gorm.io/gorm"
)

func FindUserById(id string) (models.User, *gorm.DB) {
    var user models.User

    result := database.Postgres.First(&user, id)
    return user, result
}

func FindUserByUsername(username string) (models.User, *gorm.DB) {
    var user models.User

    result := database.Postgres.Where("username = ?", username).First(&user)
    return user, result
}

func AddUser(user *models.User) (*gorm.DB) {
    return database.Postgres.Create(user)
}

