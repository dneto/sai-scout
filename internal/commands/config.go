package commands

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dneto/sai-scout/pkg/discord"
	"github.com/rs/zerolog/log"
)

var Config = func(
	saveLang func(context.Context, string, string) error,
	saveTemplate func(context.Context, string, string, string) error,
) *discord.SlashCommand {
	permissions := int64(discordgo.PermissionManageServer)
	return discord.NewCommand(&discordgo.ApplicationCommand{
		Name:                     "config",
		Description:              "Manage configuration",
		DefaultMemberPermissions: &permissions,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "language",
				Description: "Set the default language output",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "value",
						Description: "Language",
						Type:        discordgo.ApplicationCommandOptionString,
						Choices:     i18nToOptions(),
						Required:    true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "website",
				Description: "Set the website to view decks",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "template",
						Description: "Example: https://runeterra.ar/decks/code/{{code}}",
						Type:        discordgo.ApplicationCommandOptionString,
						Required:    true,
					},
					{
						Name:        "name",
						Description: "Website's name to be shown in \"View on\" button",
						Type:        discordgo.ApplicationCommandOptionString,
						Required:    true,
					},
				},
			},
		},
	}, func(s discord.Session, i *discordgo.InteractionCreate) error {

		switch i.Type {
		case discordgo.InteractionMessageComponent:
			data := i.MessageComponentData()
			split := strings.Split(data.CustomID, ";")
			switch split[1] {
			case "OK":
				label, template := split[2], split[3]
				if err := saveTemplate(context.Background(), i.GuildID, template, label); err != nil {
					log.Error().Err(err)
					return discord.ErrorResponse(s, i, err)
				}
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Flags:   discordgo.MessageFlagsEphemeral,
						Content: "Done!",
					},
				})
			case "Cancel":
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Flags:   discordgo.MessageFlagsEphemeral,
						Content: "Canceled!",
					},
				})
			}
		case discordgo.InteractionApplicationCommand:
			data := i.ApplicationCommandData()
			options := data.Options
			switch options[0].Name {
			case "language":
				language := options[0].Options[0].StringValue()
				log.Info().Str("guild", i.GuildID).Str("language", language).Msg("updating lang")
				if err := saveLang(context.Background(), i.GuildID, language); err != nil {
					log.Error().Err(err)
					return discord.ErrorResponse(s, i, err)
				}

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Flags:   discordgo.MessageFlagsEphemeral,
						Content: "Done!",
					},
				})
			case "website":
				template := options[0].Options[0].StringValue()
				label := options[0].Options[1].StringValue()

				const code = "CEDACAIFAEAQMAJJAEEAABQCA4CQCAQDAEAQWKRUAMEACAICBQBQCAIFB4AQMBJAAEDQKCQBAEAQCFA"
				templateURL := strings.Replace(template, "{{code}}", code, 1)

				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Flags: discordgo.MessageFlagsEphemeral,
					},
				})

				if err != nil {
					log.Error().Err(err)
				}

				uaaa, err := url.Parse(template)
				scheme := uaaa.Scheme
				if err != nil || (scheme != "http" && scheme != "https") {
					return discord.ErrorResponse(s, i, fmt.Errorf("Malformed URL"))
				}

				_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "Please, check if the information is correct and use the button below to view an example deck.\n" +
						"If you don't want to apply this configuration, just ignore this message.",
					Embeds: []*discordgo.MessageEmbed{
						{
							Fields: []*discordgo.MessageEmbedField{
								{Name: "template", Value: fmt.Sprintf("`%s`", template)},
								{Name: "label", Value: label},
							},
						},
					},
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									Style: discordgo.LinkButton,
									Label: fmt.Sprintf("View on %s", label),
									URL:   templateURL,
								},
							},
						},
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									Style:    discordgo.SuccessButton,
									Label:    "Apply",
									CustomID: fmt.Sprintf("config;OK;%s;%s", label, template),
								},
							},
						},
					},
				})

				if err != nil {
					log.Error().Err(err)
				}

			}
		}
		return nil
	})
}
