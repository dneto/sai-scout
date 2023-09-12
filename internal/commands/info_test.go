package commands

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/dneto/sai-scout/internal/i18n"
	"github.com/dneto/sai-scout/internal/repository"
	"github.com/dneto/sai-scout/pkg/discord"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("InfoCommand", func() {
	var localize localizeFunc = i18n.LoadTranslations().Localize
	var localizeBuildF localizeBuildFunc = localizeBuildFunc(func(language string) func(string) string {
		return func(name string) string {
			return localize(language, name)
		}
	})
	Context("InfoCommand", func() {
		var (
			command *discord.SlashCommand
		)

		BeforeEach(func() {
			command = Info(nil, nil, nil)
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

			matchName := func(ctx context.Context, language string, name string) ([]*repository.Card, error) {
				return []*repository.Card{annie}, nil
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
			resp = infoAutocompleteHandler(matchName, interaction)
		})

		It("should show all cards returned from match function", func() {
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
			handler             discord.Handler
			interactionResponse *discordgo.InteractionResponse
			followUpResp        *discordgo.WebhookParams
			err                 error
			session             discord.Session = fakeSession{
				interactionRespond: func(i *discordgo.Interaction, ir *discordgo.InteractionResponse, opts ...discordgo.RequestOption) error {
					interactionResponse = ir
					return nil
				},
				followUpMessageCreate: func(i *discordgo.Interaction, waitResponse bool, params *discordgo.WebhookParams, opts ...discordgo.RequestOption) (*discordgo.Message, error) {
					followUpResp = params
					return nil, nil
				},
			}
		)

		Context("champion card", func() {
			BeforeEach(func() {
				cardByCode := findByCodesFunc(func(ctx context.Context, language string, cardCodes ...string) ([]*repository.Card, error) {
					for _, c := range cardCodes {
						switch c {
						case "ANNIE":
							return []*repository.Card{annie}, nil
						case "ANNIELVL2":
							return []*repository.Card{annieLvl2}, nil
						default:
							return []*repository.Card{}, nil
						}
					}
					return nil, nil
				})

				handler = infoCommandHandler(cardByCode, nil, localizeBuildF)
				interaction := &discordgo.InteractionCreate{
					Interaction: &discordgo.Interaction{
						Type: discordgo.InteractionApplicationCommand,
						Data: discordgo.ApplicationCommandInteractionData{
							Options: []*discordgo.ApplicationCommandInteractionDataOption{
								{Name: "name", Value: "ANNIE"},
							},
						},
					},
				}
				err = handler(session, interaction)
			})

			It("should not return error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should send deferred message", func() {
				Expect(interactionResponse).To(Equal(&discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Processing...",
					},
				}))
			})

			It("should show embeds for all levels of the champion", func() {
				Expect(followUpResp).To(Equal(&discordgo.WebhookParams{
					Embeds: []*discordgo.MessageEmbed{
						{
							Description: "<:nx:1088276202944479342><:1_:1088291515400466524> **Annie**",
							Color:       0,
							Footer: &discordgo.MessageEmbedFooter{
								Text: "🎨 ArtistName",
							},
							Thumbnail: &discordgo.MessageEmbedThumbnail{},
							Image: &discordgo.MessageEmbedImage{
								URL: "http://path/to/annie.png",
							},
							Fields: []*discordgo.MessageEmbedField{
								{Name: "Type", Value: "Unit", Inline: true},
								{Name: "Power", Value: "0", Inline: true},
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
							Thumbnail: &discordgo.MessageEmbedThumbnail{},
							Image: &discordgo.MessageEmbedImage{
								URL:      "http://path/to/annielvl2.png",
								ProxyURL: "",
							},
							Fields: []*discordgo.MessageEmbedField{
								{Name: "Type", Value: "Unit", Inline: true},
								{Name: "Power", Value: "1", Inline: true},
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
				}))
			})
		})

		Context("spell card", func() {
			BeforeEach(func() {
				cardByCode := findByCodesFunc(func(ctx context.Context, language string, cardCodes ...string) ([]*repository.Card, error) {
					return []*repository.Card{bladesEdge}, nil
				})

				handler = infoCommandHandler(cardByCode, nil, localizeBuildF)

				interaction := &discordgo.InteractionCreate{
					Interaction: &discordgo.Interaction{
						Type: discordgo.InteractionApplicationCommand,
						Data: discordgo.ApplicationCommandInteractionData{
							Options: []*discordgo.ApplicationCommandInteractionDataOption{
								{Name: "name", Value: "BLADESEDGE"},
							},
						},
					},
				}
				err = handler(session, interaction)
			})

			It("should not return error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should send deferred message", func() {
				Expect(interactionResponse).To(Equal(&discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Processing...",
					},
				}))
			})

			It("shows card info", func() {
				Expect(followUpResp).To(Equal(&discordgo.WebhookParams{
					Embeds: []*discordgo.MessageEmbed{
						{
							Description: "<:nx:1088276202944479342><:1_:1088291515400466524> **Blade's Edge**",
							Color:       0,
							Footer: &discordgo.MessageEmbedFooter{
								Text: "🎨 ArtistName",
							},
							Thumbnail: &discordgo.MessageEmbedThumbnail{},
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
				}))
			})
		})

		Context("unit card", func() {
			BeforeEach(func() {
				cardByCode := findByCodesFunc(func(ctx context.Context, language string, cardCodes ...string) ([]*repository.Card, error) {
					return []*repository.Card{crimsonPigeon}, nil
				})

				handler = infoCommandHandler(cardByCode, nil, localizeBuildF)
				interaction := &discordgo.InteractionCreate{
					Interaction: &discordgo.Interaction{
						Type: discordgo.InteractionApplicationCommand,
						Data: discordgo.ApplicationCommandInteractionData{
							Options: []*discordgo.ApplicationCommandInteractionDataOption{
								{Name: "name", Value: "CRIMSONPIGEON"},
							},
						},
					},
				}
				err = handler(session, interaction)

			})

			It("should not return error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should send deferred message", func() {
				Expect(interactionResponse).To(Equal(&discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Processing...",
					},
				}))
			})

			It("shows card info", func() {
				Expect(followUpResp).To(Equal(&discordgo.WebhookParams{
					Embeds: []*discordgo.MessageEmbed{
						{
							Description: "<:nx:1088276202944479342><:1_:1088291515400466524> **Crimson Pigeon**",
							Color:       0,
							Footer: &discordgo.MessageEmbedFooter{
								Text: "🎨 ArtistName",
							},
							Thumbnail: &discordgo.MessageEmbedThumbnail{},
							Image: &discordgo.MessageEmbedImage{
								URL:      "http://path/to/crimsonpigeon.png",
								ProxyURL: "",
							},
							Fields: []*discordgo.MessageEmbedField{
								{Name: "Type", Value: "Unit", Inline: true},
								{Name: "Power", Value: "2", Inline: true},
								{Name: "Health", Value: "2", Inline: true},
								{Name: "Rarity", Value: "", Inline: true},
								{Name: "Description", Value: "desc", Inline: false},
								{Name: "", Value: "> _flavor_", Inline: false},
							},
						},
					},
				}))
			})
		})
	})
})
