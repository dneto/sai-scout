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

var SearchCommand = func(db Database) *discord.Command {
	return discord.NewCommand(
		&discordgo.ApplicationCommand{
			Name:        "search",
			Description: "(BETA) Show card information",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:         "name",
					Description:  "The card name",
					Type:         discordgo.ApplicationCommandOptionString,
					Required:     true,
					Autocomplete: true,
				},
			},
		},
		searchHandler(db),
	)
}

func searchHandler(db Database) func(s *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
		switch i.Type {
		case discordgo.InteractionApplicationCommandAutocomplete:
			data := i.ApplicationCommandData()

			choices := []database.Card{}

			value := data.Options[0].StringValue()
			if len(value) > 3 {
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
			}, nil

		case discordgo.InteractionApplicationCommand:
			options := i.ApplicationCommandData().Options
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
					{Name: "Keywords", Value: strings.Join(cc.Keywords, ", "), Inline: true},
					{Name: "Rarity", Value: cc.RarityRef, Inline: true},
					{Name: "Description", Value: cc.DescriptionRaw},
				}

				if cc.LevelupDescriptionRaw != "" {
					fields = append(fields,
						&discordgo.MessageEmbedField{Name: "Levelup Description", Value: cc.LevelupDescriptionRaw})
				}

				fields = append(fields,
					&discordgo.MessageEmbedField{Name: "Artist", Value: cc.ArtistName})

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
						Text: c.FlavorText,
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
