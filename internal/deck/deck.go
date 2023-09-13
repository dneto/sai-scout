package deck

import (
	"cmp"
	"context"
	"fmt"

	"github.com/dneto/sai-scout/internal/card"
	"github.com/dneto/sai-scout/internal/repository"
	"github.com/dneto/sai-scout/pkg/maps"
	"github.com/dneto/sai-scout/pkg/slices"
	lordeckcode "github.com/m0t0k1ch1/lor-deckcode-go"
	"github.com/samber/lo"
)

type DeckEntry struct {
	Count uint64
	Card  *repository.Card
}

type Deck []DeckEntry

type loadCardsInfoByCodeFunc func(ctx context.Context, language string, codes ...string) ([]*repository.Card, error)

func BuildLoadDeckInfo(loadCardsInfo loadCardsInfoByCodeFunc) func(context.Context, string, string) (Deck, error) {
	return func(ctx context.Context, language string, code string) (Deck, error) {
		deck, err := decode(code)
		if err != nil {
			return nil, err
		}

		cardsInfo, err := loadCardsInfo(ctx, language, codesFromDeck(deck)...)
		if err != nil {
			return nil, fmt.Errorf("failed to find cards: %w", err)
		}

		cardsByCode := maps.MapBy(cardsInfo, card.CardCode)

		deckWithInfo := Deck{}
		for _, de := range deck {
			if c, found := cardsByCode[de.Card.CardCode]; found {
				deckWithInfo = append(deckWithInfo, DeckEntry{
					Card:  c,
					Count: de.Count,
				})
			}
		}

		return slices.Sort(deckWithInfo, compareByCostAndName), nil
	}
}

func decode(code string) (Deck, error) {
	deck, err := lordeckcode.Decode(code)
	if err != nil {
		return nil, fmt.Errorf("failed to decode deck: %w", err)
	}

	return fromLorDeckCode(deck), nil
}

func codesFromDeck(deck Deck) []string {
	return lo.Map(deck, func(de DeckEntry, _ int) string {
		return de.Card.CardCode
	})
}

func compareByCostAndName(de DeckEntry, ee DeckEntry) int {
	c := cmp.Compare(de.Card.Cost, ee.Card.Cost)
	if c != 0 {
		return c
	}

	return cmp.Compare(de.Card.Name, ee.Card.Name)
}

func fromLorDeckCode(deck lordeckcode.Deck) Deck {
	d := make(Deck, len(deck))
	for i, ccc := range deck {
		d[i] = DeckEntry{Count: ccc.Count, Card: &repository.Card{CardCode: ccc.CardCode}}
	}

	return d
}
