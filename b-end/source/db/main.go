package db

import (
	"fmt"
	"log"
	"os"

    "nft_marketplace/eth/source/db/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Why keep Users in DB ?
// - limiting the number of listed albums per user
// - grouping listed items per user (potentially showcasing in UI)
// - grouping bids per user (showing them in UI)

func Init() {
    // Get DB variables
    dbName := os.Getenv("DB_NAME")
    userName := os.Getenv("DB_USERNAME")
    password := os.Getenv("DB_PASSWORD")
    host := os.Getenv("DB_HOST")
    env := os.Getenv("ENV")

    // Create sql info string
    sqlInfo := fmt.Sprintf("dbname=%s user=%s password=%s host=%s sslmode=disable",
                            dbName, userName, password, host); 

    // Initiate Postgresql                         
    db, err := gorm.Open(postgres.Open(sqlInfo), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    if env != "production" {
        db.AutoMigrate(&models.User{}, &models.Item{}, &models.Bid{})
    }

    fmt.Println("Successfuly connected to Postgresql DB")
}
