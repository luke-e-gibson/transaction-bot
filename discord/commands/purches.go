package command

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

func BuyRole(s *discordgo.Session) bool {
	AppCommand := discordgo.ApplicationCommand{
		Name:        "buy-role",
		Description: "Gets the users token count",
	}

	_, err := s.ApplicationCommandCreate(s.State.User.ID, "1541413756326641754", &AppCommand)
	if err != nil {
		slog.Error("Failed to register command: "+AppCommand.Name+"err: ", err)
		return false
	}

	return true

}

type BuyableRole struct {
	Price  int
	Name   string
	RoleID string
}

var (
	Roles = []*BuyableRole{
		{
			Price:  50,
			Name:   "Admin",
			RoleID: "somerole",
		},
		{
			Price:  10,
			Name:   "Somthing Else",
			RoleID: "somerolesasd",
		},
	}
)

func buyRoleHandler(s *discordgo.Session, i *discordgo.InteractionCreate, db *gorm.DB) {
	// target := getInvokingUser(i)

	options := make([]discordgo.SelectMenuOption, len(Roles))

	for i, role := range Roles {
		options[i] = discordgo.SelectMenuOption{
			Label: fmt.Sprintf("%v: $%v", role.Name, role.Price),
			Value: role.RoleID,
		}
	}

	// Create Store GUI
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pick your purchuse",
			Flags:   discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							CustomID:    "shop-selection",
							Placeholder: "Pick the role",
							Options:     options,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Buy",
							Style:    discordgo.SuccessButton,
							CustomID: "shop-buy-button",
						},
						discordgo.Button{
							Label:    "Cancel",
							Style:    discordgo.DangerButton,
							CustomID: "shop-cancel-button",
						},
					},
				},
			},
		},
	})

	if err != nil {
		slog.Error("Error while creating shop gui: ", "error", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Somthing went wrong"},
		})
	}

}

func shopCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate, db *gorm.DB) {
	switch i.ApplicationCommandData().Name {
	case "buy-role":
		buyRoleHandler(s, i, db)

	}

}

func shopMessageHandler(s *discordgo.Session, i *discordgo.InteractionCreate, db *gorm.DB) {
	slog.Info("handle me")
}

func CommandHandlerShop(s *discordgo.Session, i *discordgo.InteractionCreate, db *gorm.DB) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		shopCommandHandler(s, i, db)
	case discordgo.InteractionMessageComponent:
		shopMessageHandler(s, i, db)
	}
}
