package database

import (
	"encoding/json"
	"errors"

	"github.com/puzpuzpuz/xsync/v2"
	"github.com/sahilm/fuzzy"
)

type Database interface {
	CardByCode(code string) (Card, error)
}

type InMemory struct {
	raw    cards
	byCode *xsync.MapOf[string, *Card]
	byName *xsync.MapOf[string, *Card]
}

func NewInMemory(jsonData []byte) (*InMemory, error) {
	var cards []Card
	err := json.Unmarshal(jsonData, &cards)

	if err != nil {
		return nil, err
	}

	byCode := xsync.NewMapOf[*Card]()
	byName := xsync.NewMapOf[*Card]()

	for _, c := range cards {
		cp := &Card{}
		*cp = c
		byCode.Store(c.CardCode, cp)
		byName.Store(c.Name, cp)
	}

	inMemory := &InMemory{
		raw:    cards,
		byCode: byCode,
		byName: byName,
	}

	return inMemory, nil
}

func (i InMemory) CardByCode(code string) (Card, error) {
	card, ok := i.byCode.Load(code)

	if !ok {
		return Card{}, errors.New("not found")
	}

	return *card, nil
}

func (i InMemory) SearchByName(name string) []Card {
	results := fuzzy.FindFrom(name, i.raw)

	cards := make([]Card, len(results))
	for idx, r := range results {
		cards[idx] = i.raw[r.Index]
	}

	return cards
}
