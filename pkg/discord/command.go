package discord

import "github.com/bwmarrin/discordgo"

type Handler func(s Session, i *discordgo.InteractionCreate) error
type SlashCommand struct {
	*discordgo.ApplicationCommand
	handle Handler
}

func NewCommand(command *discordgo.ApplicationCommand, handler Handler) *SlashCommand {
	return &SlashCommand{
		ApplicationCommand: command,
		handle:             handler,
	}
}
