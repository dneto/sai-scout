package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/caarlos0/env/v9"
	"github.com/dneto/sai-scout/internal/commands"
	"github.com/dneto/sai-scout/internal/deck"
	"github.com/dneto/sai-scout/internal/i18n"
	"github.com/dneto/sai-scout/internal/repository"
	"github.com/dneto/sai-scout/pkg/discord"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/samber/mo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const lorVersion = "4.9.0"

type config struct {
	DiscordToken string `env:"DISCORD_TOKEN"`
	MongoURI     string `env:"MONGO_URI"`
}

func main() {
	cfg := config{}
	err := env.ParseWithOptions(&cfg, env.Options{RequiredIfNoDef: true})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load env vars")
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(cfg.MongoURI).SetServerAPIOptions(serverAPI)

	ctx := context.Background()
	cli, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to mongo")
	}
	defer func() {
		if err := cli.Disconnect(ctx); err != nil {
			log.Error().Err(err).Msg("Failed to close mongo connection")
		}
	}()

	err = repository.UpdateSetBundles(ctx, lorVersion, repository.InsertBuilder(cli))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to retrieve set bundles")
	}

	session, err := setupBot(cfg.DiscordToken, cli)
	defer func() {
		if err := session.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close discord connection")
		}
	}()

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to setup discord bot")
	}

	log.Info().Msg("Bot is now running.  Press CTRL-C to exit.")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, err := w.Write([]byte(""))
		if err != nil {
			fmt.Println(err)
		}
	})
	var srv http.Server
	go func() {
		err = http.ListenAndServe(":8080", nil)

		if err != nil {
			log.Fatal().Err(err).Msg("HTTP server failed")
		}
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal().Err(err).Msg("failed to shutdown http server")
	}
}

func setupBot(token string, cli *mongo.Client) (*discordgo.Session, error) {
	findCards := repository.FindCardsBuilder(cli)
	searchByName := repository.SearchByNameBuilder(cli)

	decode := deck.BuildDecode(findCards)
	localizer := i18n.LoadTranslations()

	return mo.TupleToResult(discord.NewSession(token, discordgo.IntentGuildMessages)).
		Map(discord.Open).
		Map(discord.UpdateStatus(0, fmt.Sprintf("version %s", lorVersion))).
		Map(discord.OverwriteAndHandleCommands(
			commands.Deck(decode, localizer.Localize),
			commands.Info(findCards, searchByName, localizer.Localize),
			commands.InviteCommand,
			commands.HelpCommand,
		)).
		Get()
}
