package deck

import (
	"sort"

	"github.com/dneto/sai-scout/database"
	lordeckcode "github.com/m0t0k1ch1/lor-deckcode-go"
	"github.com/samber/lo"
)

type DeckEntry struct {
	Count uint64
	Card  database.Card
}

type Deck []DeckEntry

func (d Deck) Filter(f func(de DeckEntry, _ int) bool) []DeckEntry {
	return lo.Filter(d, f)
}

func (d Deck) Champions() []DeckEntry {
	return d.Filter(func(de DeckEntry, _ int) bool {
		return de.Card.RarityRef == "Champion"
	})
}

func (d Deck) Followers() []DeckEntry {
	return d.Filter(func(de DeckEntry, _ int) bool {
		return de.Card.Type == "Unit" && de.Card.Rarity != "Champion"
	})
}

func (d Deck) Spells() []DeckEntry {
	return d.Filter(func(de DeckEntry, _ int) bool {
		return de.Card.Type == "Spell"
	})
}

func (d Deck) Landmarks() []DeckEntry {
	return d.Filter(func(de DeckEntry, _ int) bool {
		return de.Card.Type == "Landmark"
	})
}

func (d Deck) Equipments() []DeckEntry {
	return d.Filter(func(de DeckEntry, _ int) bool {
		return de.Card.Type == "Equipment"
	})
}

type Decoder struct {
	db database.Database
}

// it should return error if cannot read cards json
// it should return cardsByCode containing json (how to test this????)

func NewDecoder(db database.Database) *Decoder {

	return &Decoder{
		db: db,
	}
}

//it shouldorder cards by cost
//it should return the correct cards
//it should return error if code is invalid

func (d *Decoder) Decode(code string) (Deck, error) {
	deck, err := lordeckcode.Decode(code)
	if err != nil {
		return nil, err
	}

	sort.SliceStable(deck, func(i, j int) bool {
		ci, _ := d.db.CardByCode(deck[i].CardCode)
		cj, _ := d.db.CardByCode(deck[j].CardCode)
		return ci.Cost < cj.Cost
	})

	cardWithData := make([]DeckEntry, len(deck))
	for i, card := range deck {
		c, err := d.db.CardByCode(card.CardCode)
		if err != nil {
			continue
		}
		cardWithData[i] = DeckEntry{
			Card:  c,
			Count: card.Count,
		}
	}

	return cardWithData, nil
}
