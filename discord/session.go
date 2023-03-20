package discord

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Session struct {
	commands map[string]*Command
	session  *discordgo.Session
}

func StartSession(token string) (*Session, error) {
	discordSession, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	discordSession.Identify.Intents = discordgo.IntentGuildMessages

	err = discordSession.Open()
	if err != nil {
		return nil, err
	}

	session := &Session{
		session: discordSession,
	}

	discordSession.AddHandler(session.HandleInteraction)

	return session, nil
}

func (s *Session) Close() error {
	return s.session.Close()
}

func (s *Session) AddComand(command *Command) {
	command.Register(s.session)
	if s.commands == nil {
		s.commands = map[string]*Command{}
	}
	s.commands[command.Name] = command
}

func (s *Session) HandleInteraction(ss *discordgo.Session, i *discordgo.InteractionCreate) {
	defer func() {
		if r := recover(); r != nil {
			errorResponse(ss, i, errors.New("Unknown error"))
		}
	}()

	if c, ok := s.commands[i.ApplicationCommandData().Name]; ok {
		resp, err := c.handle(ss, i)

		if err != nil {
			err = errorResponse(ss, i, err)
			if err != nil {
				fmt.Println(err)
			}

			return
		}

		err = ss.InteractionRespond(i.Interaction, resp)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func errorResponse(s *discordgo.Session, i *discordgo.InteractionCreate, err error) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{{
				Description: fmt.Sprintf(":x: **%s**", err.Error()),
			}},
		},
	})
}
