package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dneto/sai-scout/deck"
	"github.com/dneto/sai-scout/discord"
	"github.com/dneto/sai-scout/regions"
	"github.com/samber/lo"
)

type DeckDecoder interface {
	Decode(string) (deck.Deck, error)
}

var DeckCommand = func(decoder DeckDecoder) *discord.Command {
	return discord.NewCommand(
		&discordgo.ApplicationCommand{
			Name:        "deck",
			Description: "Show cards present in given deck",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "code",
					Description: "Legends of Runeterra deck code",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
		handler(decoder),
	)
}

func handler(decoder DeckDecoder) func(s *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
		options := i.ApplicationCommandData().Options
		deckCode := options[0].Value.(string)

		de, err := decoder.Decode(deckCode)

		if err != nil {
			return nil, errors.New("invalid code")
		}

		cardsByType := map[string][]deck.DeckEntry{
			"Champions":  de.Champions(),
			"Followers":  de.Followers(),
			"Spells":     de.Spells(),
			"Landmarks":  de.Landmarks(),
			"Equipments": de.Equipments(),
		}

		typesShowOrder := []string{"Champions", "Followers", "Spells", "Landmarks", "Equipments"}
		fields := make([]*discordgo.MessageEmbedField, len(typesShowOrder))
		for i, t := range typesShowOrder {
			cardsByTypeT := cardsByType[t]

			if len(cardsByTypeT) == 0 {
				continue
			}

			cards := lo.Map(cardsByTypeT, func(de deck.DeckEntry, _ int) string {
				return cardToStr(de)
			})

			fields[i] = &discordgo.MessageEmbedField{
				Name:   t,
				Value:  strings.Join(cards, "\n"),
				Inline: true,
			}
		}

		response := &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:  deckCode,
					Fields: fields,
				},
			},
		}

		if i.Interaction.Member != nil {
			icon := fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", i.Interaction.Member.User.ID, i.Interaction.Member.User.Avatar)
			response.Embeds[0].Footer = &discordgo.MessageEmbedFooter{
				Text:    i.Interaction.Member.User.Username,
				IconURL: icon,
			}
		}

		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: response,
		}, nil
	}
}

func cardToStr(c deck.DeckEntry) string {
	rs := ""
	for _, r := range c.Card.RegionRefs {
		rs = rs + regions.Emote(r)
	}

	return fmt.Sprintf("%s %s %s | **%d**", costEmoji[c.Card.Cost], rs, c.Card.Name, c.Count)
}

var costEmoji = map[int]string{
	0:  "<:0_:1086337952310906983>",
	1:  "<:1_:1086331931286851654>",
	2:  "<:2_:1086333429311877170>",
	3:  "<:3_:1086333430712782948>",
	4:  "<:4_:1086333433086746816>",
	5:  "<:5_:1086333434433122504>",
	6:  "<:6_:1086333436811296899>",
	7:  "<:7_:1086333437876650050>",
	8:  "<:8_:1086333439407554592>",
	9:  "<:9_:1086333441785737366>",
	10: "<:10:1086333443085971508>",
	11: "<:11:1086333445350891741>",
	12: "<:12:1086333481979756584>",
	13: "<:13:1086333447645188106>",
	14: "<:14:1086333450295984228>",
	15: "<:15:1086333483032518768>",
}
