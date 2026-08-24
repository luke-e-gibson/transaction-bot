package Env

import (
	"os"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	DiscordBotToken string
	DatabaseUrl     string
}

func LoadEnv() (EnvConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return EnvConfig{}, err
	}

	return EnvConfig{
		DiscordBotToken: os.Getenv("DISCORD_BOT_TOKEN"),
		DatabaseUrl:     os.Getenv("DATABASE_URL"),
	}, nil
}
