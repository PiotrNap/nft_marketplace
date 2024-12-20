package models

import (
	"time"
)

type Bid struct {
    ID              uint      `gorm:"primaryKey"`
    IsActive        bool      `gorm:"not null"`
    UserID          uint      `gorm:"not null"`
    User            User      `gorm:"foreignKey:UserID"`
    Price           uint      `gorm:"not null"`
    TxHash          string    `gorm:"size:66;not null"` 
}

type Item struct {
    ID              uint        `gorm:"primaryKey"`
    IsActive        bool        `gorm:"not null"`
    CreationDate    time.Time   `gorm:"not null"`
    ContractAddress string      `gorm:"size:44; not null"`
    TokenId         uint        `gorm:"not null"`
    MinPrice        uint        `gorm:"not null"`
    Duration        time.Time   `gorm:"not null"`
    CurrentBids     []Bid       `gorm:"foreignKey:UserID; not null"`
    Owner           User        `gorm:"foreignKey:OwnerID"`
    OwnerID         uint        `gorm:"not null"`
}

type User struct {
    ID              uint      `gorm:"primaryKey"`
    Username        string    `gorm:"not null"`
    Challenge       string    `gorm:"size:66"`
    Salt            string    `gorm:"size:36"`
    Verified        bool      `gorm:"not null; default=false"`
    PubKey          string    `gorm:"size:66; index"`
    Address         string    `gorm:"size:42; not null; index"`
    Items           []Item    `gorm:"foreignKey:OwnerID"`
    Bids            []Bid     `gorm:"foreignKey:UserID"`
}
