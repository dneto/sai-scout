package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dneto/sai-scout/commands"
	"github.com/dneto/sai-scout/database"
	"github.com/dneto/sai-scout/deck"
	"github.com/dneto/sai-scout/discord"
)

func main() {
	token := os.Getenv("DISCORD_TOKEN")

	if token == "" {
		log.Fatal("enviroment variable DISCORD_TOKEN is empty")
	}

	cardsJSON := os.Getenv("CARDS_FILE")
	if cardsJSON == "" {
		cardsJSON = "cards.json"
	}

	s, err := discord.StartSession(token)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	data, err := os.ReadFile(cardsJSON)
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewInMemory(data)
	if err != nil {
		log.Fatal(err)
	}

	dec := deck.NewDecoder(db)

	s.AddComand(commands.DeckCommand(dec))
	s.AddComand(commands.InfoCommand(db))
	err = s.CleanCommands()
	if err != nil {
		log.Println("Error while remove old commands", err)
	}

	log.Print("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
