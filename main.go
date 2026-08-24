package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	DB "transaction-bot/db"
	"transaction-bot/discord"
	Env "transaction-bot/env"
)

func main() {
	slog.Info("Starting...")
	env, err := Env.LoadEnv()
	if err != nil {
		slog.Error("Could not Parse Env")
		return
	}

	slog.Info("Connecting to database")
	db, err := DB.NewDb(env)
	if err != nil {
		slog.Error("Error Connecting To Database")
		return
	}

	slog.Info("Migrareing database")
	if err := DB.MigrareDB(db); err != nil {
		slog.Error("Database Magration Failed, Error: ", err)
		return
	}

	slog.Info("Starting discord bot")
	session := discord.StartBot(env, db)

	defer session.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}
