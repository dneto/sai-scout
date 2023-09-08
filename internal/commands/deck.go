package commands

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dneto/sai-scout/internal/card"
	"github.com/dneto/sai-scout/internal/deck"
	"github.com/dneto/sai-scout/internal/i18n"
	"github.com/dneto/sai-scout/internal/regions"
	"github.com/dneto/sai-scout/internal/repository"
	"github.com/dneto/sai-scout/pkg/discord"
	"github.com/samber/lo"
	"golang.org/x/text/width"
)

type DeckDecoder interface {
	Decode(i18n.Locale, string) (deck.Deck, error)
}

var Deck = func(decoder decodeFunc, localize localizeFunc) *discord.SlashCommand {
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
				{
					Name:        "language",
					Description: "Language",
					Type:        discordgo.ApplicationCommandOptionString,
					Choices:     i18nToOptions(),
					Required:    false,
				},
			},
		},
		deckCommandHandler(decoder, localize),
	)
}

func i18nToOptions() []*discordgo.ApplicationCommandOptionChoice {
	opts := make([]*discordgo.ApplicationCommandOptionChoice, len(i18n.Locales))
	for i, locale := range i18n.Locales {
		opts[i] = &discordgo.ApplicationCommandOptionChoice{
			Name:  locale.Name(),
			Value: string(locale),
		}
	}
	return opts
}

func deckCommandHandler(decode decodeFunc, localize localizeFunc) func(s discord.Session, i *discordgo.InteractionCreate) error {
	return func(s discord.Session, i *discordgo.InteractionCreate) error {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title: "**Processing**",
					},
				},
			},
		})

		i.ApplicationCommandData()
		options := i.ApplicationCommandData().Options
		deckCode := options[0].Value.(string)
		language := string(i18n.Default)
		if len(options) > 1 {
			language = options[1].StringValue()
		}

		decodedDeck, err := decode(context.Background(), language, deckCode)

		if err != nil {
			fmt.Println(err)
			return discord.ErrorResponse(s, i, errors.New("invalid code"))
		}

		filter := func(d deck.Deck, predicate func(c *repository.Card) bool) []deck.DeckEntry {
			return lo.Filter(d, func(den deck.DeckEntry, _ int) bool {
				return predicate(den.Card)
			})
		}

		cardsByType := map[string][]deck.DeckEntry{
			"Champions":  filter(decodedDeck, card.IsChampion),
			"Followers":  filter(decodedDeck, card.IsFollower),
			"Spells":     filter(decodedDeck, card.IsSpell),
			"Landmarks":  filter(decodedDeck, card.IsLandmark),
			"Equipments": filter(decodedDeck, card.IsEquipment),
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
				Name:   localize(language, t),
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

		embeds := []*discordgo.MessageEmbed{
			{
				Title:  deckCode,
				Fields: fields,
			},
		}

		if i.Interaction.Member != nil {
			name := i.Interaction.Member.Nick
			if name == "" {
				name = i.Interaction.Member.User.Username
			}

			icon := avatarURL(i)
			embeds[0].Footer = &discordgo.MessageEmbedFooter{
				Text:    name,
				IconURL: icon,
			}
		}

		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Embeds: embeds,
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
		})

		return err
	}
}

func avatarURL(i *discordgo.InteractionCreate) string {
	member := i.Interaction.Member
	avatar := member.Avatar
	if avatar != "" {
		return fmt.Sprintf("https://cdn.discordapp.com/guilds/%s/users/%s/avatars/%s.png", i.GuildID, member.User.ID, avatar)
	}

	avatar = member.User.Avatar
	return fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", member.User.ID, avatar)
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
