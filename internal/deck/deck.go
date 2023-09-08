package deck

import (
	"cmp"
	"context"
	"fmt"
	"slices"

	lordeckcode "github.com/dneto/lor-deckcode-go"
	"github.com/dneto/sai-scout/internal/repository"
	"github.com/samber/lo"
	"github.com/samber/mo"
)

type DeckEntry struct {
	Count uint64
	Card  *repository.Card
}

type Deck []DeckEntry

type findByCodesFunc func(ctx context.Context, language string, codes ...string) ([]*repository.Card, error)

func BuildDecode(findByCodes findByCodesFunc) func(context.Context, string, string) (Deck, error) {
	return func(ctx context.Context, language string, code string) (Deck, error) {
		return mo.TupleToResult(decode(code)).
			Map(populate(ctx, language, findByCodes)).Get()
	}
}

func decode(code string) (Deck, error) {
	toDeck := func(ccc lordeckcode.CardCodeAndCount, _ int) DeckEntry {
		return DeckEntry{Count: ccc.Count, Card: &repository.Card{CardCode: ccc.CardCode}}
	}

	deck, err := lordeckcode.Decode(code)
	if err != nil {
		return nil, fmt.Errorf("failed to decode deck: %w", err)
	}

	return lo.Map(deck, toDeck), nil
}

func populate(ctx context.Context, locale string, findByCodes findByCodesFunc) func(Deck) (Deck, error) {
	return func(d Deck) (Deck, error) {
		codes := lo.Map(d, func(de DeckEntry, _ int) string {
			return de.Card.CardCode
		})

		cards, err := findByCodes(ctx, locale, codes...)
		if err != nil {
			return nil, fmt.Errorf("failed to find cards: %w", err)
		}

		loadCard := func(de DeckEntry, _ int) DeckEntry {
			for _, c := range cards {
				if c.CardCode == de.Card.CardCode {
					return DeckEntry{
						Card:  c,
						Count: de.Count,
					}
				}
			}
			return DeckEntry{}
		}

		dd := lo.Map(d, loadCard)
		slices.SortFunc(dd, compareDeckEntry)
		return dd, nil
	}
}

func compareDeckEntry(de DeckEntry, ee DeckEntry) int {
	c := cmp.Compare(de.Card.Cost, ee.Card.Cost)
	if c != 0 {
		return c
	}

	return cmp.Compare(de.Card.Name, ee.Card.Name)
}
