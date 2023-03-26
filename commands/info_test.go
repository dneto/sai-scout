package commands

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/dneto/sai-scout/database"
	"github.com/dneto/sai-scout/discord"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("InfoCommand", func() {
	Context("InfoCommand", func() {
		var (
			command *discord.Command
		)

		BeforeEach(func() {
			command = InfoCommand(nil)
		})

		It("should have name 'info'", func() {
			Expect(command.ApplicationCommand.Name).To(Equal("info"))
		})

		It("should have option 'name'", func() {
			Expect(command.ApplicationCommand.Options[0].Name).To(Equal("name"))
		})
	})

	Context("autocomplete", func() {
		var (
			resp *discordgo.InteractionResponse
		)
		BeforeEach(func() {
			db := fakeDB{}

			db.fakeSearchByName = func(name string) []database.Card {
				return []database.Card{annie, annieLvl2}
			}

			interaction := &discordgo.InteractionCreate{
				Interaction: &discordgo.Interaction{
					Type: discordgo.InteractionApplicationCommandAutocomplete,
					Data: discordgo.ApplicationCommandInteractionData{
						Options: []*discordgo.ApplicationCommandInteractionDataOption{
							{
								Type:  discordgo.ApplicationCommandOptionString,
								Value: "Annie",
							},
						},
					},
				},
			}
			resp = infoAutocompleteHandler(db, interaction)
		})

		It("should show embeds for all levels of the champion", func() {
			Expect(resp).To(Equal(&discordgo.InteractionResponse{
				Type: discordgo.InteractionApplicationCommandAutocompleteResult,
				Data: &discordgo.InteractionResponseData{
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{Name: "Annie", Value: "ANNIE"},
					},
				},
			}))
		})
	})

	Context("handler", func() {
		var (
			resp *discordgo.InteractionResponse
			err  error
		)
		Context("champion card", func() {
			BeforeEach(func() {
				db := fakeDB{}
				db.fakeCardByCode = func(cardCode string) (database.Card, error) {
					switch cardCode {
					case "ANNIE":
						return annie, nil
					case "ANNIELVL2":
						return annieLvl2, nil
					default:
						return database.Card{}, nil
					}
				}

				h := infoCommandHandler(db)
				session := &discordgo.Session{}
				interaction := &discordgo.InteractionCreate{
					Interaction: &discordgo.Interaction{
						Type: discordgo.InteractionApplicationCommand,
						Data: discordgo.ApplicationCommandInteractionData{
							Options: []*discordgo.ApplicationCommandInteractionDataOption{
								{Value: "ANNIE"},
							},
						},
					},
				}
				resp, err = h(session, interaction)
			})

			It("returns no error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should show embeds for all levels of the champion", func() {
				Expect(resp).To(Equal(&discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							{
								Description: "<:nx:1088276202944479342><:1_:1088291515400466524> **Annie**",
								Color:       0,
								Footer: &discordgo.MessageEmbedFooter{
									Text: "🎨 ArtistName",
								},
								Image: &discordgo.MessageEmbedImage{
									URL: "http://path/to/annie.png",
								},
								Fields: []*discordgo.MessageEmbedField{
									{Name: "Type", Value: "Unit", Inline: true},
									{Name: "Attack", Value: "0", Inline: true},
									{Name: "Health", Value: "2", Inline: true},
									{Name: "Rarity", Value: "Champion", Inline: true},
									{Name: "Description", Value: "desc annie", Inline: false},
									{Name: "Level Up", Value: "lvl up desc", Inline: false},
									{
										Name:   "",
										Value:  "> _flavor annie_",
										Inline: false,
									},
								},
							},
							{
								Description: "<:nx:1088276202944479342><:1_:1088291515400466524> **Annie**",
								Footer: &discordgo.MessageEmbedFooter{
									Text: "🎨 ArtistName",
								},
								Image: &discordgo.MessageEmbedImage{
									URL:      "http://path/to/annielvl2.png",
									ProxyURL: "",
								},
								Fields: []*discordgo.MessageEmbedField{
									{Name: "Type", Value: "Unit", Inline: true},
									{Name: "Attack", Value: "1", Inline: true},
									{Name: "Health", Value: "2", Inline: true},
									{Name: "Description", Value: "desc annie lvl 2", Inline: false},
									{
										Name:   "",
										Value:  "> _flavor annie lvl2_",
										Inline: false,
									},
								},
							},
						},
					},
				}))
			})
		})

		Context("spell card", func() {
			BeforeEach(func() {
				db := fakeDB{}
				db.fakeCardByCode = func(cardCode string) (database.Card, error) {
					return bladesEdge, nil
				}

				h := infoCommandHandler(db)
				session := &discordgo.Session{}
				interaction := &discordgo.InteractionCreate{
					Interaction: &discordgo.Interaction{
						Type: discordgo.InteractionApplicationCommand,
						Data: discordgo.ApplicationCommandInteractionData{
							Options: []*discordgo.ApplicationCommandInteractionDataOption{
								{Value: "BLADESEDGE"},
							},
						},
					},
				}
				resp, err = h(session, interaction)
			})

			It("returns no error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("shows card info", func() {
				Expect(resp).To(Equal(&discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							{
								Description: "<:nx:1088276202944479342><:1_:1088291515400466524> **Blade's Edge**",
								Color:       0,
								Footer: &discordgo.MessageEmbedFooter{
									Text: "🎨 ArtistName",
								},
								Image: &discordgo.MessageEmbedImage{
									URL:      "http://path/to/bladesedge.png",
									ProxyURL: "",
								},
								Fields: []*discordgo.MessageEmbedField{
									{Name: "Type", Value: "Spell", Inline: true},
									{Name: "Keywords", Value: "Fast", Inline: true},
									{Name: "Rarity", Value: "Common", Inline: true},
									{Name: "Description", Value: "desc", Inline: false},
									{
										Name:   "",
										Value:  "> _flavor_",
										Inline: false,
									},
								},
							},
						},
					},
				}))
			})
		})

		Context("unit card", func() {
			BeforeEach(func() {
				db := fakeDB{}
				db.fakeCardByCode = func(cardCode string) (database.Card, error) {
					return crimsonPigeon, nil
				}

				h := infoCommandHandler(db)
				session := &discordgo.Session{}
				interaction := &discordgo.InteractionCreate{
					Interaction: &discordgo.Interaction{
						Type: discordgo.InteractionApplicationCommand,
						Data: discordgo.ApplicationCommandInteractionData{
							Options: []*discordgo.ApplicationCommandInteractionDataOption{
								{Value: "CRIMSONPIGEON"},
							},
						},
					},
				}
				resp, err = h(session, interaction)
			})

			It("returns no error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("shows card info", func() {
				Expect(resp).To(Equal(&discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							{
								Description: "<:nx:1088276202944479342><:1_:1088291515400466524> **Crimson Pigeon**",
								Color:       0,
								Footer: &discordgo.MessageEmbedFooter{
									Text: "🎨 ArtistName",
								},
								Image: &discordgo.MessageEmbedImage{
									URL:      "http://path/to/crimsonpigeon.png",
									ProxyURL: "",
								},
								Fields: []*discordgo.MessageEmbedField{
									{Name: "Type", Value: "Unit", Inline: true},
									{Name: "Attack", Value: "2", Inline: true},
									{Name: "Health", Value: "2", Inline: true},
									{Name: "Rarity", Value: "", Inline: true},
									{Name: "Description", Value: "desc", Inline: false},
									{Name: "", Value: "> _flavor_", Inline: false},
								},
							},
						},
					},
				}))
			})
		})
	})
})

type fakeDB struct {
	fakeCardByCode   func(cardCode string) (database.Card, error)
	fakeSearchByName func(name string) []database.Card
}

func (f fakeDB) CardByCode(cardCode string) (database.Card, error) {
	if f.fakeCardByCode != nil {
		return f.fakeCardByCode(cardCode)
	}

	return database.Card{}, errors.New("not implemented")
}

func (f fakeDB) SearchByName(name string) []database.Card {
	if f.fakeSearchByName != nil {
		return f.fakeSearchByName(name)
	}

	return nil
}
