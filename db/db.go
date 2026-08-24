package DB

import (
	"fmt"
	"log/slog"
	"time"
	Env "transaction-bot/env"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDb(env Env.EnvConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(env.DatabaseUrl), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		slog.Error("Failed to connect to database")
		return nil, fmt.Errorf("Failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("Failed to open to database")
		return nil, err
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	slog.Info("Connected to database")

	return db, nil
}

func MigrareDB(db *gorm.DB) error {
	return db.AutoMigrate(&UserEntry{}, &TransactionEntry{})
}
