package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dneto/sai-scout/discord"
)

var InviteCommand = discord.NewCommand(
	&discordgo.ApplicationCommand{
		Name:        "invite",
		Description: "Send invite link",
	},
	inviteCommandHandler,
)

func inviteCommandHandler(s *discordgo.Session, in *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Style: discordgo.LinkButton,
							Label: "Invite",
							URL:   "https://discord.com/api/oauth2/authorize?client_id=1086224659231559680&permissions=0&scope=bot",
						},
					},
				},
			},
		},
	}, nil
}
