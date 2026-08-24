package discord

import (
	"log/slog"
	command "transaction-bot/discord/commands"
	Env "transaction-bot/env"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

var (
	integerOptionMinValue          = 1.0
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionManageGuild
)

func StartBot(env Env.EnvConfig, db *gorm.DB) *discordgo.Session {
	s, err := discordgo.New("Bot " + env.DiscordBotToken)
	if err != nil {
		slog.Error("Failed to create discord bot")
		return nil
	}

	// Register Events
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		slog.Info("Logged in as: " + s.State.User.Username + "#" + s.State.User.Discriminator)
	})
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		command.CommandHandlerEconomy(s, i, db)
		command.CommandHandlerShop(s, i, db)

	})

	err = s.Open()
	if err != nil {
		slog.Error("Failed to Create Bot Session")
		return nil
	}

	// Register Commands
	slog.Info("Registering Commands")
	command.UpVoteResponse(s)
	command.GetUserTokens(s)
	command.BuyRole(s)

	return s
}
