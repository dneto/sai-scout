package deck_test

import (
	"context"
	"errors"
	"testing"

	"github.com/dneto/sai-scout/internal/deck"
	"github.com/dneto/sai-scout/internal/i18n"
	"github.com/dneto/sai-scout/internal/repository"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDeck(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Decks Suite")
}

var _ = Describe("Decks", func() {
	var ctx context.Context = context.Background()
	Context("BuildDecode", func() {
		var (
			target deck.Deck
			err    error

			findByCodeFunc func(ctx context.Context, language string, codes ...string) ([]*repository.Card, error)
		)

		Context("no errors", func() {
			BeforeEach(func() {
				findByCodeFunc = func(_ context.Context, _ string, codes ...string) ([]*repository.Card, error) {
					cards := make([]*repository.Card, 0)
					for _, code := range codes {
						card, found := cardByCode[code]
						if !found {
							return nil, errors.New("not found")
						}
						cards = append(cards, card)
					}
					return cards, nil
				}
				deckCode := "CEAAAAICAYBQYHA"
				target, err = deck.BuildDecode(findByCodeFunc)(ctx, string(i18n.Default), deckCode)

			})

			It("returns no error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("cards should be sorted by cost", func() {
				Expect(target).To(Equal(deck.Deck{
					{Count: 1, Card: annie},
					{Count: 1, Card: ravenbloomConservatory},
				}))
			})
		})

		Context("find errors", func() {
			BeforeEach(func() {
				findByCodeFunc = func(_ context.Context, _ string, codes ...string) ([]*repository.Card, error) {
					return nil, errors.New("fatal error")
				}
				deckCode := "CEAAAAICAYBQYHA"
				target, err = deck.BuildDecode(findByCodeFunc)(ctx, string(i18n.Default), deckCode)

			})

			It("returns no error", func() {
				Expect(err).To(HaveOccurred())
			})

			It("should match err", func() {
				Expect(err).To(MatchError("failed to find cards: fatal error"))
			})
		})

		Context("decode errors", func() {
			BeforeEach(func() {
				findByCodeFunc = func(_ context.Context, _ string, codes ...string) ([]*repository.Card, error) {
					return nil, errors.New("fatal error")
				}
				deckCode := ""
				target, err = deck.BuildDecode(findByCodeFunc)(ctx, string(i18n.Default), deckCode)

			})

			It("returns no error", func() {
				Expect(err).To(HaveOccurred())
			})

			It("should match err", func() {
				Expect(err).To(MatchError("failed to decode deck: failed to read format and version: EOF"))
			})
		})
	})
})

var (
	annie                  = &repository.Card{CardCode: "06NX012", Name: "Annie", Cost: 0, RarityRef: "Champion"}
	ravenbloomConservatory = &repository.Card{CardCode: "06NX028", Name: "Ravenbloom Conservatory", Cost: 1, Type: "Landmark"}
	theDarkinBallista      = &repository.Card{CardCode: "06NX020", Name: "The Darkin Ballista", Cost: 2, Type: "Equipment"}
	crimsonPigeon          = &repository.Card{CardCode: "06NX041", Name: "Crimson Pigeon", Cost: 3, Type: "Unit"}
	bladesEdge             = &repository.Card{CardCode: "01NX043", Name: "Blade's Edge", Cost: 4, Type: "Spell"}

	cardByCode map[string]*repository.Card = map[string]*repository.Card{
		"06NX012": annie,
		"06NX028": ravenbloomConservatory,
		"06NX020": theDarkinBallista,
		"06NX041": crimsonPigeon,
		"01NX043": bladesEdge,
	}
)
