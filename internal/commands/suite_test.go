package commands

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/dneto/sai-scout/internal/repository"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBooks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commands Suite")
}

var (
	annie = &repository.Card{
		CardCode:              "ANNIE",
		AssociatedCardRefs:    []string{"ANNIELVL2"},
		Name:                  "Annie",
		Cost:                  1,
		Rarity:                "Champion",
		RarityRef:             "Champion",
		RegionRefs:            []string{"Noxus"},
		TypeRef:               "Unit",
		Type:                  "Unit",
		Attack:                0,
		Health:                2,
		DescriptionRaw:        "desc annie",
		FlavorText:            "flavor annie",
		LevelupDescriptionRaw: "lvl up desc",
		Assets: []struct {
			GameAbsolutePath string "json:\"gameAbsolutePath\""
			FullAbsolutePath string "json:\"fullAbsolutePath\""
		}{
			{FullAbsolutePath: "http://path/to/annie.png"},
		},
		SupertypeRef: "Champion",
		ArtistName:   "ArtistName",
	}

	annieLvl2 = &repository.Card{
		CardCode:       "ANNIELVL2",
		Name:           "Annie",
		Cost:           1,
		Rarity:         "None",
		RarityRef:      "None",
		RegionRefs:     []string{"Noxus"},
		Type:           "Unit",
		TypeRef:        "Unit",
		Attack:         1,
		Health:         2,
		DescriptionRaw: "desc annie lvl 2",
		FlavorText:     "flavor annie lvl2",
		Assets: []struct {
			GameAbsolutePath string "json:\"gameAbsolutePath\""
			FullAbsolutePath string "json:\"fullAbsolutePath\""
		}{
			{FullAbsolutePath: "http://path/to/annielvl2.png"},
		},
		Supertype:    "Champion",
		SupertypeRef: "Champion",
		ArtistName:   "ArtistName",
	}
	ravenbloomConservatory = &repository.Card{Name: "Ravenbloom Conservatory", Cost: 1, Type: "Landmark", RegionRefs: []string{"Noxus"}, TypeRef: "Landmark"}
	theDarkinBallista      = &repository.Card{Name: "The Darkin Ballista", Cost: 1, Type: "Equipment", RegionRefs: []string{"Noxus"}, TypeRef: "Equipment"}

	crimsonPigeon = &repository.Card{
		Name:           "Crimson Pigeon",
		Cost:           1,
		Type:           "Unit",
		TypeRef:        "Unit",
		RegionRefs:     []string{"Noxus"},
		Attack:         2,
		Health:         2,
		DescriptionRaw: "desc",
		FlavorText:     "flavor",
		Assets: []struct {
			GameAbsolutePath string "json:\"gameAbsolutePath\""
			FullAbsolutePath string "json:\"fullAbsolutePath\""
		}{
			{FullAbsolutePath: "http://path/to/crimsonpigeon.png"},
		},
		ArtistName: "ArtistName",
	}

	bladesEdge = &repository.Card{
		Name:           "Blade's Edge",
		Cost:           1,
		Type:           "Spell",
		TypeRef:        "Spell",
		RegionRefs:     []string{"Noxus"},
		Rarity:         "Common",
		DescriptionRaw: "desc",
		FlavorText:     "flavor",
		Assets: []struct {
			GameAbsolutePath string "json:\"gameAbsolutePath\""
			FullAbsolutePath string "json:\"fullAbsolutePath\""
		}{
			{FullAbsolutePath: "http://path/to/bladesedge.png"},
		},
		Keywords:   []string{"Fast"},
		Supertype:  "Champion",
		ArtistName: "ArtistName",
	}
)

type fakeSession struct {
	followUpMessageCreate func(i *discordgo.Interaction, waitResponse bool, params *discordgo.WebhookParams, opts ...discordgo.RequestOption) (*discordgo.Message, error)
	interactionRespond    func(i *discordgo.Interaction, ir *discordgo.InteractionResponse, opts ...discordgo.RequestOption) error
}

func (s fakeSession) FollowupMessageCreate(i *discordgo.Interaction, waitResponse bool, params *discordgo.WebhookParams, opts ...discordgo.RequestOption) (*discordgo.Message, error) {
	return s.followUpMessageCreate(i, waitResponse, params, opts...)
}

func (s fakeSession) InteractionRespond(i *discordgo.Interaction, ir *discordgo.InteractionResponse, opts ...discordgo.RequestOption) error {
	return s.interactionRespond(i, ir, opts...)
}
