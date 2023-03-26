package commands

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/dneto/sai-scout/deck"
	"github.com/dneto/sai-scout/discord"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeckCommand", func() {
	Context("DeckCommand", func() {
		var (
			command *discord.Command
		)

		BeforeEach(func() {
			command = DeckCommand(nil)
		})

		It("should have name 'deck'", func() {
			Expect(command.ApplicationCommand.Name).To(Equal("deck"))
		})

		It("should have option 'code'", func() {
			Expect(command.ApplicationCommand.Options[0].Name).To(Equal("code"))
		})
	})

	Context("handler", func() {
		var (
			resp *discordgo.InteractionResponse
			err  error
		)

		Context("sucess", func() {
			BeforeEach(func() {
				decoder := fakeDeckDecoder{}
				decoder.fakeDecode = func(deckCode string) (deck.Deck, error) {
					return deck.Deck{
						deck.DeckEntry{Card: annie, Count: 1},
						deck.DeckEntry{Card: ravenbloomConservatory, Count: 1},
						deck.DeckEntry{Card: crimsonPigeon, Count: 1},
						deck.DeckEntry{Card: theDarkinBallista, Count: 1},
						deck.DeckEntry{Card: bladesEdge, Count: 1},
					}, nil
				}

				h := deckCommandHandler(decoder)
				session := &discordgo.Session{}
				interaction := &discordgo.InteractionCreate{
					Interaction: &discordgo.Interaction{
						Type: discordgo.InteractionApplicationCommand,
						Data: discordgo.ApplicationCommandInteractionData{
							Options: []*discordgo.ApplicationCommandInteractionDataOption{
								{Value: "DECKCODE"},
							},
						},
					},
				}
				resp, err = h(session, interaction)

			})

			It("returns no error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should have a embed", func() {
				Expect(resp.Data.Embeds).ToNot(BeEmpty())
			})

			It("embed title is the deck code", func() {
				Expect(resp.Data.Embeds[0].Title).To(Equal("DECKCODE"))
			})

			It("should have a embed field for each card type", func() {
				Expect(resp.Data.Embeds[0].Fields).To(Equal(
					[]*discordgo.MessageEmbedField{
						{
							Name:   "Champions",
							Value:  "**１** <:nx:1088276202944479342><:1_:1088291515400466524> Annie",
							Inline: true,
						},
						{
							Name:   "Followers",
							Value:  "**１** <:nx:1088276202944479342><:1_:1088291515400466524> Crimson Pigeon",
							Inline: true,
						},
						{
							Name:   "Spells",
							Value:  "**１** <:nx:1088276202944479342><:1_:1088291515400466524> Blade's Edge",
							Inline: true,
						},
						{
							Name:   "Landmarks",
							Value:  "**１** <:nx:1088276202944479342><:1_:1088291515400466524> Ravenbloom Conservatory",
							Inline: true,
						},
						{
							Name:   "Equipments",
							Value:  "**１** <:nx:1088276202944479342><:1_:1088291515400466524> The Darkin Ballista",
							Inline: true,
						},
					},
				))
			})

			It("should have Runeterra AR button", func() {
				Expect(resp.Data.Components[0]).To(Equal(
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								Style: discordgo.LinkButton,
								Label: "View on Runeterra AR",
								URL:   "https://runeterra.ar/decks/code/" + "DECKCODE",
							},
						},
					},
				))
			})
		})

		Context("contains more than 10 cards of a type", func() {
			BeforeEach(func() {
				decoder := fakeDeckDecoder{}
				decoder.fakeDecode = func(deckCode string) (deck.Deck, error) {
					return deck.Deck{
						deck.DeckEntry{Card: annie, Count: 1},
						deck.DeckEntry{Card: annie, Count: 1},
						deck.DeckEntry{Card: annie, Count: 1},
						deck.DeckEntry{Card: annie, Count: 1},
						deck.DeckEntry{Card: annie, Count: 1},
						deck.DeckEntry{Card: annie, Count: 1},
						deck.DeckEntry{Card: annie, Count: 1},
						deck.DeckEntry{Card: annie, Count: 1},
						deck.DeckEntry{Card: annie, Count: 1},
						deck.DeckEntry{Card: annie, Count: 1},
						deck.DeckEntry{Card: annie, Count: 1},
					}, nil
				}

				h := deckCommandHandler(decoder)
				session := &discordgo.Session{}
				interaction := &discordgo.InteractionCreate{
					Interaction: &discordgo.Interaction{
						Type: discordgo.InteractionApplicationCommand,
						Data: discordgo.ApplicationCommandInteractionData{
							Options: []*discordgo.ApplicationCommandInteractionDataOption{
								{Value: "DECKCODE"},
							},
						},
					},
				}
				resp, err = h(session, interaction)
			})

			It("should have a embed field for each card type", func() {
				Expect(resp.Data.Embeds[0].Fields).To(Equal(
					[]*discordgo.MessageEmbedField{
						{
							Name:   "Champions",
							Value:  "**１** <:nx:1088276202944479342><:1_:1088291515400466524> Annie\n**１** <:nx:1088276202944479342><:1_:1088291515400466524> Annie\n**１** <:nx:1088276202944479342><:1_:1088291515400466524> Annie\n**１** <:nx:1088276202944479342><:1_:1088291515400466524> Annie\n**１** <:nx:1088276202944479342><:1_:1088291515400466524> Annie\n**１** <:nx:1088276202944479342><:1_:1088291515400466524> Annie\n**１** <:nx:1088276202944479342><:1_:1088291515400466524> Annie\n**１** <:nx:1088276202944479342><:1_:1088291515400466524> Annie\n**１** <:nx:1088276202944479342><:1_:1088291515400466524> Annie\n**１** <:nx:1088276202944479342><:1_:1088291515400466524> Annie",
							Inline: true,
						},
						{
							Name:   "ㅤ",
							Value:  "**１** <:nx:1088276202944479342><:1_:1088291515400466524> Annie",
							Inline: true,
						},
					},
				))
			})
		})

		Context("fields without cards are not showed", func() {

			BeforeEach(func() {
				decoder := fakeDeckDecoder{}
				decoder.fakeDecode = func(deckCode string) (deck.Deck, error) {
					return deck.Deck{
						deck.DeckEntry{Card: annie, Count: 1},
					}, nil
				}

				h := deckCommandHandler(decoder)
				session := &discordgo.Session{}
				interaction := &discordgo.InteractionCreate{
					Interaction: &discordgo.Interaction{
						Type: discordgo.InteractionApplicationCommand,
						Data: discordgo.ApplicationCommandInteractionData{
							Options: []*discordgo.ApplicationCommandInteractionDataOption{
								{Value: "DECKCODE"},
							},
						},
					},
				}
				resp, err = h(session, interaction)
			})

			It("should have only one field", func() {
				Expect(resp.Data.Embeds[0].Fields).To(HaveLen(1))
			})

			It("field should be champion field", func() {
				Expect(resp.Data.Embeds[0].Fields[0].Name).To(Equal("Champions"))
			})
		})

		Context("invalid deck code", func() {
			BeforeEach(func() {
				decoder := fakeDeckDecoder{}
				decoder.fakeDecode = func(deckCode string) (deck.Deck, error) {
					return nil, errors.New("error while decoding")
				}

				h := deckCommandHandler(decoder)
				session := &discordgo.Session{}
				interaction := &discordgo.InteractionCreate{
					Interaction: &discordgo.Interaction{
						Type: discordgo.InteractionApplicationCommand,
						Data: discordgo.ApplicationCommandInteractionData{
							Options: []*discordgo.ApplicationCommandInteractionDataOption{
								{Value: "DECKCODE"},
							},
						},
					},
				}
				resp, err = h(session, interaction)
			})

			It("returns invalid deck code error", func() {
				Expect(err).To(MatchError("invalid code"))
			})
		})
	})
})

type fakeDeckDecoder struct {
	fakeDecode func(deckCode string) (deck.Deck, error)
}

func (f fakeDeckDecoder) Decode(cardCode string) (deck.Deck, error) {
	if f.fakeDecode != nil {
		return f.fakeDecode(cardCode)
	}

	return nil, errors.New("not implemented")
}
