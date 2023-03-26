package commands

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dneto/sai-scout/database"
	"github.com/dneto/sai-scout/discord"
	"github.com/dneto/sai-scout/regions"
	"github.com/samber/lo"
)

type Database interface {
	SearchByName(string) []database.Card
	CardByCode(string) (database.Card, error)
}

var InfoCommand = func(db Database) *discord.Command {
	return discord.NewCommand(
		&discordgo.ApplicationCommand{
			Name:        "info",
			Description: "Show card information like regions, cost, name, keywords, description, etc.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:         "name",
					Description:  "The card name (autocomplete)",
					Type:         discordgo.ApplicationCommandOptionString,
					Required:     true,
					Autocomplete: true,
				},
			},
		},
		infoCommandHandler(db),
	)
}

func infoCommandHandler(db Database) func(s *discordgo.Session, in *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
	return func(s *discordgo.Session, in *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
		switch in.Type {
		case discordgo.InteractionApplicationCommandAutocomplete:
			return infoAutocompleteHandler(db, in), nil
		case discordgo.InteractionApplicationCommand:
			options := in.ApplicationCommandData().Options
			cardName := options[0].Value.(string)

			c, err := db.CardByCode(cardName)

			if err != nil {
				return nil, errors.New("invalid code")
			}

			cards := []database.Card{}
			if c.Supertype == "Champion" && c.Type == "Unit" {
				associated := c.AssociatedCardRefs
				for _, a := range associated {
					cc, _ := db.CardByCode(a)
					if cc.Supertype == "Champion" && cc.Type == "Unit" {
						cards = append(cards, cc)
					}
				}
			}

			cards = append(cards, c)

			sort.SliceStable(cards, func(i, j int) bool {
				return cards[i].CardCode < cards[j].CardCode
			})

			embeds := []*discordgo.MessageEmbed{}
			for _, cc := range cards {
				fields := []*discordgo.MessageEmbedField{
					{Name: "Type", Value: cc.Type, Inline: true},
				}

				if cc.Type == "Unit" {
					fields = append(fields,
						&discordgo.MessageEmbedField{Name: "Attack", Value: fmt.Sprint(cc.Attack), Inline: true})

					fields = append(fields,
						&discordgo.MessageEmbedField{Name: "Health", Value: fmt.Sprint(cc.Health), Inline: true})

				}

				if len(cc.Keywords) > 0 {
					fields = append(fields,
						&discordgo.MessageEmbedField{Name: "Keywords", Value: strings.Join(cc.Keywords, "\n"), Inline: true})

				}

				if cc.RarityRef != "None" {
					fields = append(fields,
						&discordgo.MessageEmbedField{Name: "Rarity", Value: cc.RarityRef, Inline: true})
				}

				if cc.DescriptionRaw != "" {
					fields = append(fields,
						&discordgo.MessageEmbedField{Name: "Description", Value: cc.DescriptionRaw})
				}

				if cc.LevelupDescriptionRaw != "" {
					fields = append(fields,
						&discordgo.MessageEmbedField{Name: "Level Up", Value: cc.LevelupDescriptionRaw})
				}

				flavorLines := strings.Split(cc.FlavorText, "\n")
				flavorLines = lo.Map(flavorLines, func(line string, _ int) string {
					return fmt.Sprintf("> _%s_", line)
				})

				fields = append(fields,
					&discordgo.MessageEmbedField{Name: "", Value: strings.Join(flavorLines, "\n")})

				regs := ""
				for _, r := range cc.RegionRefs {
					regs = regs + regions.Emote(r)
				}

				embed := &discordgo.MessageEmbed{
					Description: fmt.Sprintf("%s%s **%s**", regs, costEmoji[cc.Cost], cc.Name),
					Image: &discordgo.MessageEmbedImage{
						URL: cc.Assets[0].FullAbsolutePath,
					},
					Fields: fields,
					Footer: &discordgo.MessageEmbedFooter{
						Text: "🎨 " + c.ArtistName,
					},
				}
				embeds = append(embeds, embed)
			}

			response := &discordgo.InteractionResponseData{
				Embeds: embeds,
			}

			return &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: response,
			}, nil
		}

		return nil, nil
	}
}

func infoAutocompleteHandler(db Database, in *discordgo.InteractionCreate) *discordgo.InteractionResponse {
	data := in.ApplicationCommandData()

	value := data.Options[0].StringValue()
	choices := []database.Card{}
	if len(value) > 1 {
		choices = db.SearchByName(value)
	}

	choices = lo.UniqBy(choices, func(c database.Card) string {
		return c.Name
	})

	if len(choices) > 25 {
		choices = choices[:25]
	}

	ch := make([]*discordgo.ApplicationCommandOptionChoice, len(choices))
	for i, c := range choices {
		name := c.Name
		ch[i] = &discordgo.ApplicationCommandOptionChoice{
			Name:  name,
			Value: c.CardCode,
		}
	}

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: ch,
		},
	}
}
