package DB

import "gorm.io/gorm"

type UserEntry struct {
	ID                string `gorm:"primaryKey"` // Discord ID
	DisplayName       string
	CurrentTokenCount int64 `gorm:"default:0"`
}

type TransactionEntry struct {
	gorm.Model
	UserFromID *string   `gorm:"index"` // nil = system-issued
	UserFrom   UserEntry `gorm:"foreignKey:UserFromID;references:ID"`
	UserToID   string    `gorm:"index"`
	UserTo     UserEntry `gorm:"foreignKey:UserToID;references:ID"`
	Amount     int64
	Reason     string // optional: "upvote", "daily bonus", "admin grant", etc.
}
