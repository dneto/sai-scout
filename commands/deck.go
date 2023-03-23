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
	"golang.org/x/text/width"
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
		deckCommandHandler(decoder),
	)
}

func deckCommandHandler(decoder DeckDecoder) func(s *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
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
		fields := []*discordgo.MessageEmbedField{}
		for _, t := range typesShowOrder {
			if len(cardsByType[t]) == 0 {
				continue
			}

			cards := lo.Map(cardsByType[t], func(de deck.DeckEntry, _ int) string {
				return cardToStr(de)
			})

			cs := lo.Chunk(cards, 10)

			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   t,
				Value:  strings.Join(cs[0], "\n"),
				Inline: true,
			})

			for _, c := range cs[1:] {
				title := "ㅤ"
				inline := true
				fields = append(fields, &discordgo.MessageEmbedField{
					Name:   title,
					Value:  strings.Join(c, "\n"),
					Inline: inline,
				})
			}
		}

		response := &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:  deckCode,
					Fields: fields,
				},
			},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Style: discordgo.LinkButton,
							Label: "View on Runeterra AR",
							URL:   "https://runeterra.ar/decks/code/" + deckCode,
						},
					},
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

	n := width.Widen.String(fmt.Sprint(c.Count))
	return fmt.Sprintf("**%s** %s%s %s", n, rs, costEmoji[c.Card.Cost], c.Card.Name)
}

var costEmoji = map[int]string{
	0:  "<:0_:1088290496906018869>",
	1:  "<:1_:1088291515400466524>",
	2:  "<:2_:1088291517069787147>",
	3:  "<:3_:1088291520072925274>",
	4:  "<:4_:1088291522589495406>",
	5:  "<:5_:1088291525303210034>",
	6:  "<:6_:1088291527899480084>",
	7:  "<:7_:1088293024435556394>",
	8:  "<:8_:1088293026541097020>",
	9:  "<:9_:1088293029070262282>",
	10: "<:10:1088293030471155762>",
	11: "<:11:1088293032845115464>",
	12: "<:12:1088293034787086346>",
	13: "<:13:1088293037261717524>",
	14: "<:14:1088293040155807755>",
	15: "<:15:1088293041594433637>",
	16: "<:16:1088293042957602816>",
	17: "<:17:1088293045490958356>",
}
