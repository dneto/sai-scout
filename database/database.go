package database

import (
	"encoding/json"
	"errors"

	"github.com/puzpuzpuz/xsync/v2"
)

type Database interface {
	CardByCode(code string) (Card, error)
}

type InMemory struct {
	m *xsync.MapOf[string, Card]
}

func NewInMemory(jsonData []byte) (*InMemory, error) {
	var cards []Card
	err := json.Unmarshal(jsonData, &cards)

	if err != nil {
		return nil, err
	}

	m := xsync.NewMapOf[Card]()
	for _, c := range cards {
		m.Store(c.CardCode, c)
	}

	inMemory := &InMemory{
		m: m,
	}

	return inMemory, nil
}

func (i InMemory) CardByCode(code string) (Card, error) {
	card, ok := i.m.Load(code)

	if !ok {
		return Card{}, errors.New("not found")
	}

	return card, nil
}
