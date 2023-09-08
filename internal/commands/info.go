package commands

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/bwmarrin/discordgo"
	"github.com/dneto/sai-scout/internal/i18n"
	"github.com/dneto/sai-scout/internal/regions"
	"github.com/dneto/sai-scout/internal/repository"
	"github.com/dneto/sai-scout/pkg/discord"
	"github.com/dneto/sai-scout/pkg/discord/embed"
	"github.com/dneto/sai-scout/pkg/discord/option"
	"github.com/samber/lo"
)

var infoCommand = &discordgo.ApplicationCommand{
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
		{
			Name:        "language",
			Description: "Language",
			Type:        discordgo.ApplicationCommandOptionString,
			Choices:     i18nToOptions(),
			Required:    false,
		},
	},
}

var Info = func(findByCodes findByCodesFunc, matchName matchNameFunc, localize localizeFunc) *discord.SlashCommand {
	lfunc := func(l string) func(string) string {
		return func(s string) string {
			return localize(l, s)
		}
	}
	return discord.NewCommand(infoCommand, infoCommandHandler(findByCodes, matchName, lfunc))
}

func infoCommandHandler(
	findByCodes findByCodesFunc,
	matchName matchNameFunc,
	localizeBuilder localizeBuildFunc) discord.Handler {

	return func(s discord.Session, in *discordgo.InteractionCreate) error {
		switch in.Type {
		case discordgo.InteractionApplicationCommandAutocomplete:
			return s.InteractionRespond(in.Interaction, infoAutocompleteHandler(matchName, in))
		case discordgo.InteractionApplicationCommand:
			s.InteractionRespond(in.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Processing...",
				},
			})

			options := in.ApplicationCommandData().Options
			cardName, err := option.Get[string](options, "name")
			if err != nil {
				return discord.ErrorResponse(s, in, errors.New("missing name"))
			}

			language := option.GetOrElse(options, "language", string(i18n.Default))
			localize := localizeBuilder(language)

			ctx := context.Background()
			cs, err := findByCodes(ctx, language, cardName)
			if err != nil || len(cs) < 1 {
				return discord.ErrorResponse(s, in, errors.New("invalid code: "+cardName))
			}

			c := cs[0]
			cards := []*repository.Card{}
			if c.SupertypeRef == "Champion" && c.TypeRef == "Unit" {
				associated, err := findByCodes(ctx, language, c.AssociatedCardRefs...)
				if err != nil {
					log.Println(err)
				}

				for _, a := range associated {
					if a.SupertypeRef == "Champion" && a.TypeRef == "Unit" {
						cards = append(cards, a)
					}
				}
			}

			cards = append(cards, c)

			sort.SliceStable(cards, func(i, j int) bool {
				return cards[i].CardCode < cards[j].CardCode
			})

			embeds := lo.Map(cards, cardToEmbed(localize))
			s.FollowupMessageCreate(in.Interaction, false, &discordgo.WebhookParams{
				Embeds: embeds,
			})

			return nil
		}

		return nil
	}
}

func cardToEmbed(localize func(string) string) func(c *repository.Card, _ int) *discordgo.MessageEmbed {
	return func(c *repository.Card, _ int) *discordgo.MessageEmbed {
		me := &discordgo.MessageEmbed{
			Description: buildTitle(c),
			Image: &discordgo.MessageEmbedImage{
				URL: c.Assets[0].FullAbsolutePath,
			},
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: c.Assets[0].GameAbsolutePath,
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: "🎨 " + c.ArtistName,
			},
		}

		addFields := embed.AddFields(me)

		addFields(embed.InlineField(localize("Type"), c.Type))

		if c.TypeRef == "Unit" || c.TypeRef == "Equipment" {
			addFields(
				embed.InlineField(localize("Attack"), strconv.Itoa(c.Attack)),
				embed.InlineField(localize("Health"), strconv.Itoa(c.Health)),
			)
		}

		if len(c.Keywords) > 0 {
			addFields(embed.InlineField(localize("Keywords"), strings.Join(c.Keywords, "\n")))
		}

		if c.RarityRef != "None" {
			addFields(embed.InlineField(localize("Rarity"), c.Rarity))
		}

		if c.DescriptionRaw != "" {
			addFields(embed.Field(localize("Description"), c.DescriptionRaw))
		}

		if c.LevelupDescriptionRaw != "" {
			addFields(
				embed.Field(localize("Level Up"), c.LevelupDescriptionRaw),
			)
		}

		flavorLines := strings.Split(c.FlavorText, "\n")
		flavorLines = lo.Map(flavorLines, func(line string, _ int) string {
			return fmt.Sprintf("> _%s_", line)
		})

		addFields(embed.Field("", strings.Join(flavorLines, "\n")))

		if len(c.Formats) > 0 && c.Collectible {
			addFields(embed.InlineField(localize("Formats"), strings.Join(c.Formats, ", ")))
		}

		return me
	}
}

func infoAutocompleteHandler(matchName matchNameFunc, in *discordgo.InteractionCreate) *discordgo.InteractionResponse {
	data := in.ApplicationCommandData()
	o := data.Options

	name, _ := option.Get[string](o, "name")
	language := option.GetOrElse(o, "language", string(i18n.Default))

	choices, err := matchName(context.Background(), language, name)
	if err != nil {
		log.Println(err)
	}

	ch := make([]*discordgo.ApplicationCommandOptionChoice, len(choices))
	for i, c := range choices {
		ch[i] = &discordgo.ApplicationCommandOptionChoice{
			Name:  c.Name,
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

func buildTitle(c *repository.Card) string {
	cardTemplate := "{{.regions}}{{.cost}} **{{.name}}**"
	tmpl, err := template.New("").Parse(cardTemplate)
	if err != nil {
		return ""
	}

	regs := ""
	for _, r := range c.RegionRefs {
		regs = regs + regions.Emote(r)
	}
	bbuf := bytes.NewBuffer([]byte{})
	err = tmpl.Execute(bbuf, map[string]string{
		"regions": regs,
		"cost":    costEmoji[c.Cost],
		"name":    c.Name,
	})

	if err != nil {
		return ""
	}

	return bbuf.String()
}
