package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dneto/sai-scout/discord"
)

var HelpCommand = discord.NewCommand(
	&discordgo.ApplicationCommand{
		Name:        "help",
		Description: "Show help text for command",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "deck",
				Description: "Show help for deck command",
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "info",
				Description: "Show help for info command",
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "invite",
				Description: "Show help for invite command",
			},
		},
	},
	helpCommandHandler,
)

func helpCommandHandler(s *discordgo.Session, in *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
	data := in.ApplicationCommandData()
	options := data.Options
	switch options[0].Name {
	case "deck":
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
				Content: "> **/deck** `code`: Shows cards of the given deck code and a button for Runeterra AR view" + "\n" +
					"> " + "\n" +
					"> _Options:_" + "\n" +
					"> • `code`: Legends of Runeterra deck code. Example: CICACAIDFABAOBQ2DMCAEBQUCYWTUBIHAMAQMBYIBEBACBYGDEAQOAYDAA",
			},
		}, nil

	case "info":
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
				Content: "> **/info** `name`: Shows a card info" + "\n" +
					"> " + "\n" +
					"> _Options:_" + "\n" +
					"> • `name`: The card name. As soon as you start typing, the field will filter and show cards name matching what you typed ",
			},
		}, nil

	case "invite":
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "> **/invite**: Send bot invite link",
			},
		}, nil

	default:
		return nil, nil
	}
}
