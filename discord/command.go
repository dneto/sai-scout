package discord

import "github.com/bwmarrin/discordgo"

type handler func(s *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error)

type Command struct {
	*discordgo.ApplicationCommand
	handle handler
}

func NewCommand(command *discordgo.ApplicationCommand, handler handler) *Command {
	return &Command{
		ApplicationCommand: command,
		handle:             handler,
	}
}

func (c Command) Register(s *discordgo.Session) error {
	commands, err := s.ApplicationCommands(s.State.User.ID, "")
	if err != nil {
		return err
	}

	for _, co := range commands {
		if c.Name == co.Name {
			_, err := s.ApplicationCommandEdit(s.State.User.ID, "", co.ID, c.ApplicationCommand)
			return err
		}
	}

	_, err = s.ApplicationCommandCreate(s.State.User.ID, "", c.ApplicationCommand)
	return err
}
