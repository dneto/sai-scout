package embed

import "github.com/bwmarrin/discordgo"

func Field(name string, value string) *discordgo.MessageEmbedField {
	return &discordgo.MessageEmbedField{
		Name:  name,
		Value: value,
	}
}

func InlineField(name string, value string) *discordgo.MessageEmbedField {
	field := Field(name, value)
	field.Inline = true
	return field
}

func AddFields(me *discordgo.MessageEmbed) func(...*discordgo.MessageEmbedField) {
	return func(mef ...*discordgo.MessageEmbedField) {
		me.Fields = append(me.Fields, mef...)
	}
}
