package discord

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

type Session interface {
	InteractionRespond(i *discordgo.Interaction, ir *discordgo.InteractionResponse, opts ...discordgo.RequestOption) error
	FollowupMessageCreate(i *discordgo.Interaction, waitResponse bool, params *discordgo.WebhookParams, opts ...discordgo.RequestOption) (*discordgo.Message, error)
}

func NewSession(token string, intents discordgo.Intent) (*discordgo.Session, error) {
	session, err := discordgo.New("Bot " + token)

	if err != nil {
		return nil, err
	}

	session.Identify.Intents = intents
	return session, nil
}

func Open(s *discordgo.Session) (*discordgo.Session, error) {
	return s, s.Open()
}

func Close(s *discordgo.Session) error {
	return s.Close()
}

func UpdateStatus(idle int, msg string) func(*discordgo.Session) (*discordgo.Session, error) {
	return func(s *discordgo.Session) (*discordgo.Session, error) {
		return s, s.UpdateWatchStatus(0, msg)
	}
}

func OverwriteAndHandleCommands(commands ...*SlashCommand) func(*discordgo.Session) (*discordgo.Session, error) {
	return func(s *discordgo.Session) (*discordgo.Session, error) {
		commandMap := make(map[string]*SlashCommand)
		for _, c := range commands {
			commandMap[c.Name] = c
		}
		err := OverwriteCommands(s, commandMap)
		if err != nil {
			return s, fmt.Errorf("failed to overwrite commands")
		}

		s.AddHandler(HandleInteractionCreate(commandMap))
		return s, nil
	}
}

func OverwriteCommands(s *discordgo.Session, newCommands map[string]*SlashCommand) error {
	appID := s.State.User.ID
	oldCommands, err := s.ApplicationCommands(appID, "")
	if err != nil {
		return err
	}

	discordCommands := make(map[string]*discordgo.ApplicationCommand)
	for _, c := range oldCommands {
		discordCommands[c.Name] = c
	}

	for name, oldCommand := range discordCommands {
		if _, ok := newCommands[name]; !ok {
			log.Debug().Str("command", name).Msg("deleting command")
			err := s.ApplicationCommandDelete(s.State.User.ID, "", oldCommand.ID)
			if err != nil {
				log.Error().Err(err).Str("command", name).
					Msg("failed to delete application command")
			}
			continue
		}
	}

	for name, newCommand := range newCommands {
		if oldCommand, ok := discordCommands[name]; ok {
			log.Debug().Str("command", name).Msg("editing command")
			_, err := s.ApplicationCommandEdit(appID, "", oldCommand.ID, newCommand.ApplicationCommand)
			if err != nil {
				log.Error().Err(err).Str("command", name).Msg("failed to edit application command")
				return err
			}
			continue
		}

		log.Debug().Str("command", name).Msg("creating command")
		_, err := s.ApplicationCommandCreate(appID, "", newCommand.ApplicationCommand)
		if err != nil {
			log.Error().Err(err).Str("command", name).Msg("failed to create application command")
			return err
		}
	}

	return nil
}

func HandleInteractionCreate(commands map[string]*SlashCommand) func(*discordgo.Session, *discordgo.InteractionCreate) {
	return func(ss *discordgo.Session, ic *discordgo.InteractionCreate) {
		defer func() {
			if r := recover(); r != nil {
				log.Error().Any("recover", r).Msg("")
				ErrorResponse(ss, ic, errors.New("Unknown error"))
			}
		}()
		switch ic.Type {
		case discordgo.InteractionApplicationCommand:
			if c, ok := commands[ic.ApplicationCommandData().Name]; ok {
				err := c.handle(ss, ic)

				if err != nil {
					log.Error().Err(err).Msg("Failed do respond interaction")
				}
			}
		case discordgo.InteractionMessageComponent:
			customID := ic.MessageComponentData().CustomID
			split := strings.Split(customID, ";")
			if c, ok := commands[split[0]]; ok {
				err := c.handle(ss, ic)

				if err != nil {
					log.Error().Err(err).Msg("Failed do respond interaction")
				}
			}
		}
	}
}

func ErrorResponse(s Session, i *discordgo.InteractionCreate, err error) error {
	log.Error().Err(err).Msg("Failed do respond interaction")

	_, err = s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
		Content: "### <:poroshock:1151475832494247947> **Error**",
		Embeds: []*discordgo.MessageEmbed{
			{
				Description: err.Error(),
			},
		},
	})
	return err
}
