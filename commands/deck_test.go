package commands

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/dneto/sai-scout/database"
	"github.com/dneto/sai-scout/deck"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeckCommand", func() {
	Context("handler", func() {
		var (
			resp *discordgo.InteractionResponse
			err  error
		)

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

			h := handler(decoder)
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
						Value:  "<:1_:1086331931286851654> <:nx:1086330895214071921> Annie | **1**",
						Inline: true,
					},
					{
						Name:   "Followers",
						Value:  "<:1_:1086331931286851654> <:nx:1086330895214071921> Crimson Pigeon | **1**",
						Inline: true,
					},
					{
						Name:   "Spells",
						Value:  "<:1_:1086331931286851654> <:nx:1086330895214071921> Blade's Edge | **1**",
						Inline: true,
					},
					{
						Name:   "Landmarks",
						Value:  "<:1_:1086331931286851654> <:nx:1086330895214071921> Ravenbloom Conservatory | **1**",
						Inline: true,
					},
					{
						Name:   "Equipments",
						Value:  "<:1_:1086331931286851654> <:nx:1086330895214071921> The Darkin Ballista | **1**",
						Inline: true,
					},
				},
			))
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

var (
	annie                  = database.Card{Name: "Annie", Cost: 1, RarityRef: "Champion", RegionRefs: []string{"Noxus"}}
	ravenbloomConservatory = database.Card{Name: "Ravenbloom Conservatory", Cost: 1, Type: "Landmark", RegionRefs: []string{"Noxus"}}
	theDarkinBallista      = database.Card{Name: "The Darkin Ballista", Cost: 1, Type: "Equipment", RegionRefs: []string{"Noxus"}}
	crimsonPigeon          = database.Card{Name: "Crimson Pigeon", Cost: 1, Type: "Unit", RegionRefs: []string{"Noxus"}}
	bladesEdge             = database.Card{Name: "Blade's Edge", Cost: 1, Type: "Spell", RegionRefs: []string{"Noxus"}}
)
