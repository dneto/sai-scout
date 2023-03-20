package deck

import (
	"errors"

	"github.com/dneto/sai-scout/database"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Decoder", func() {
	Context("success", func() {
		var (
			deck Deck
			err  error
		)

		BeforeEach(func() {
			fakedb := fakeDB{}
			fakedb.fakeCardByCode = func(cardCode string) (database.Card, error) {
				return map[string]database.Card{
					"01PZ008": teemo,
					"01PZ054": boomcrewRookie,
				}[cardCode], nil
			}
			decoder := NewDecoder(fakedb)
			deck, err = decoder.Decode("CEAAAAICAECAQNQ")

		})

		It("returns no error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("cards should be sorted by cost", func() {
			Expect(deck[0].Card.Cost).To(BeNumerically("<", deck[1].Card.Cost))
		})

		It("return the cards from database", func() {
			Expect(deck[0].Card).To(Equal(teemo))
		})
	})

	Context("cannot find card", func() {
		var (
			deck Deck
		)

		BeforeEach(func() {
			fakedb := fakeDB{}
			fakedb.fakeCardByCode = func(cardCode string) (database.Card, error) {
				return database.Card{}, errors.New("not found")
			}
			decoder := NewDecoder(fakedb)
			deck, _ = decoder.Decode("CEAAAAICAECAQNQ")

		})

		It("returns deck with missing card as empty", func() {
			Expect(deck).To(Equal(Deck{{}, {}}))
		})

	})
})

var _ = Describe("Deck", func() {
	Context("Deck", func() {
		var deck Deck
		var err error

		BeforeEach(func() {
			fakedb := fakeDB{}
			fakedb.fakeCardByCode = func(cardCode string) (database.Card, error) {
				return map[string]database.Card{
					"06NX012": annie,
					"06NX028": ravenbloomConservatory,
					"06NX020": theDarkinBallista,
					"06NX041": crimsonPigeon,
					"01NX043": bladesEdge,
				}[cardCode], nil
			}
			decoder := NewDecoder(fakedb)
			deck, err = decoder.Decode("CEAAAAQBAEBSWBAGAMGBIHBJ")
		})

		It("should not return error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return champions", func() {
			Expect(deck.Champions()).To(Equal(
				[]DeckEntry{{Count: 1, Card: annie}},
			))
		})

		It("should return followers", func() {
			Expect(deck.Followers()).To(Equal(
				[]DeckEntry{{Count: 1, Card: crimsonPigeon}},
			))
		})

		It("should return spells", func() {
			Expect(deck.Spells()).To(Equal(
				[]DeckEntry{{Count: 1, Card: bladesEdge}},
			))
		})

		It("should return equipments", func() {
			Expect(deck.Equipments()).To(Equal(
				[]DeckEntry{{Count: 1, Card: theDarkinBallista}},
			))
		})

		It("should return landmarks", func() {
			Expect(deck.Landmarks()).To(Equal(
				[]DeckEntry{{Count: 1, Card: ravenbloomConservatory}},
			))

		})
	})
})

type fakeDB struct {
	fakeCardByCode func(cardCode string) (database.Card, error)
}

func (f fakeDB) CardByCode(cardCode string) (database.Card, error) {
	if f.fakeCardByCode != nil {
		return f.fakeCardByCode(cardCode)
	}

	return database.Card{}, errors.New("not implemented")
}

var (
	teemo                  = database.Card{Name: "Teemo", Cost: 1, RarityRef: "Champion"}
	boomcrewRookie         = database.Card{Name: "BoomcrewRookie", Cost: 2, Type: "Unit"}
	annie                  = database.Card{Name: "Annie", Cost: 1, RarityRef: "Champion"}
	ravenbloomConservatory = database.Card{Name: "Ravenbloom Conservatory", Cost: 1, Type: "Landmark"}
	theDarkinBallista      = database.Card{Name: "The Darkin Ballista", Cost: 1, Type: "Equipment"}
	crimsonPigeon          = database.Card{Name: "Crimson Pigeon", Cost: 1, Type: "Unit"}
	bladesEdge             = database.Card{Name: "Blade's Edge", Cost: 1, Type: "Spell"}
)
