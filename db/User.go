package DB

import (
	"context"
	"errors"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func EnsureUserInDatabase(dUser *discordgo.User, db *gorm.DB) (*UserEntry, error) {
	var user UserEntry

	err := db.Where(UserEntry{ID: dUser.ID}).Attrs(UserEntry{
		DisplayName:       dUser.DisplayName(),
		CurrentTokenCount: 0,
	}).FirstOrCreate(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUser(dUser *discordgo.User, db *gorm.DB) (*UserEntry, error) {
	res, err := gorm.G[UserEntry](db).Where(UserEntry{ID: dUser.ID}).First(context.TODO())
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func CreateSystemTransaction(db *gorm.DB, toID string, amount int64, reason string) (*TransactionEntry, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	var txn TransactionEntry

	err := db.Transaction(func(tx *gorm.DB) error {
		var to UserEntry
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&to, "id = ?", toID).Error; err != nil {
			return err
		}

		if err := tx.Model(&to).
			Update("current_token_count", gorm.Expr("current_token_count + ?", amount)).Error; err != nil {
			return err
		}

		txn = TransactionEntry{
			UserFromID: nil, // system-issued
			UserToID:   toID,
			Amount:     amount,
			Reason:     reason,
		}
		if err := tx.Create(&txn).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &txn, nil
}
