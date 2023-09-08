package option

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func Get[T any](options []*discordgo.ApplicationCommandInteractionDataOption, name string) (T, error) {
	var option T
	for _, o := range options {
		if o.Name == name {
			option, ok := o.Value.(T)
			if !ok {
				return option, fmt.Errorf("unable to convert %s to %v", o.Name, option)
			}
			return option, nil
		}
	}

	return option, fmt.Errorf("not found")
}

func GetOrElse[T any](options []*discordgo.ApplicationCommandInteractionDataOption, name string, def T) T {
	t, err := Get[T](options, name)
	if err != nil {
		return def
	}
	return t
}
