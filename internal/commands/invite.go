package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dneto/sai-scout/pkg/discord"
)

var InviteCommand = discord.NewCommand(
	&discordgo.ApplicationCommand{
		Name:        "invite",
		Description: "Send invite link",
	},
	inviteCommandHandler,
)

func inviteCommandHandler(s discord.Session, in *discordgo.InteractionCreate) error {
	return s.InteractionRespond(in.Interaction, &discordgo.InteractionResponse{
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
	})
}
