package commands

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/dneto/sai-scout/internal/deck"
	"github.com/dneto/sai-scout/internal/repository"
)

type findByCodesFunc func(ctx context.Context, language string, cardCodes ...string) ([]*repository.Card, error)
type matchNameFunc func(ctx context.Context, language string, name string) ([]*repository.Card, error)
type decodeFunc func(ctx context.Context, language string, code string) (deck.Deck, error)
type localizeFunc func(language string, messageID string) string
type localizeBuildFunc func(string) func(string) string

type interactionHandler func(s *discordgo.Session, in *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error)
