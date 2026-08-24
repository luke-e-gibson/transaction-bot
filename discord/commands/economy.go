package command

import (
	"fmt"
	"log/slog"
	DB "transaction-bot/db"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

func UpVoteResponse(s *discordgo.Session) bool {
	AppCommand := discordgo.ApplicationCommand{
		Name:        "upvote",
		Description: "Give User Tokens",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "user",
				Description: "The User To Give The Token",
				Type:        discordgo.ApplicationCommandOptionUser,
			},
		},
	}

	_, err := s.ApplicationCommandCreate(s.State.User.ID, "1541413756326641754", &AppCommand)
	if err != nil {
		slog.Error("Failed to register command: "+AppCommand.Name+"err: ", err)
		return false
	}

	return true
}

func GetUserTokens(s *discordgo.Session) bool {
	AppCommand := discordgo.ApplicationCommand{
		Name:        "get-tokens",
		Description: "Gets the users token count",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "user",
				Description: "The user to list otherwise your self",
				Type:        discordgo.ApplicationCommandOptionUser,
			},
		},
	}

	_, err := s.ApplicationCommandCreate(s.State.User.ID, "1541413756326641754", &AppCommand)
	if err != nil {
		slog.Error("Failed to register command: "+AppCommand.Name+"err: ", err)
		return false
	}

	return true

}

func getInvokingUser(i *discordgo.InteractionCreate) *discordgo.User {
	if i.Member != nil {
		return i.Member.User // guild context
	}
	return i.User // DM context
}

func upVoteResponseHandler(s *discordgo.Session, i *discordgo.InteractionCreate, db *gorm.DB) {
	data := i.ApplicationCommandData()

	var target *discordgo.User
	for _, opt := range data.Options {
		if opt.Name == "user" {
			target = opt.UserValue(s)
		}
	}

	if target == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "You must specify a user."},
		})
		return
	}

	invoker := getInvokingUser(i)

	_, err := DB.EnsureUserInDatabase(invoker, db)
	if err != nil {
		slog.Info("Could Not Create Invoking User Entry")

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Somthing went wrong."}})

		return
	}

	targetUser, terr := DB.EnsureUserInDatabase(target, db)
	if terr != nil {
		slog.Info("Could Not Create Target User Entry")

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Somthing went wrong."}})

		return
	}

	transaction, tterr := DB.CreateSystemTransaction(db, targetUser.ID, 10, "upvote")
	if tterr != nil {
		slog.Info("Could Not Create Transaction")
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Somthing went wrong."}})

		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: fmt.Sprintf("Gave 10 tokens to <@%s> (txnb #%d)", targetUser.ID, transaction.ID)},
	})
}

func getUserTokensHandler(s *discordgo.Session, i *discordgo.InteractionCreate, db *gorm.DB) {
	data := i.ApplicationCommandData()

	var target *discordgo.User
	for _, opt := range data.Options {
		if opt.Name == "user" {
			target = opt.UserValue(s)
		}
	}

	if target == nil {
		target = getInvokingUser(i)
	}

	if target == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Somthing went wrong"},
		})
		return
	}

	user, err := DB.GetUser(target, db)
	if err != nil {
		slog.Error("Failed to user")
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Somthing went wrong"},
		})

		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: fmt.Sprintf("<@%s> has %d tokens to there name", target.ID, user.CurrentTokenCount)},
	})
}

func CommandHandlerEconomy(s *discordgo.Session, i *discordgo.InteractionCreate, db *gorm.DB) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	switch i.ApplicationCommandData().Name {
	case "upvote":
		upVoteResponseHandler(s, i, db)

	case "get-tokens":
		getUserTokensHandler(s, i, db)
	}

}
