package commands

import (
	"testing"

	"github.com/dneto/sai-scout/database"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBooks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commands Suite")
}

var (
	annie = database.Card{
		CardCode:              "ANNIE",
		AssociatedCardRefs:    []string{"ANNIELVL2"},
		Name:                  "Annie",
		Cost:                  1,
		RarityRef:             "Champion",
		RegionRefs:            []string{"Noxus"},
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
		Supertype:  "Champion",
		ArtistName: "ArtistName",
	}

	annieLvl2 = database.Card{
		CardCode:       "ANNIELVL2",
		Name:           "Annie",
		Cost:           1,
		RarityRef:      "None",
		RegionRefs:     []string{"Noxus"},
		Type:           "Unit",
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
		Supertype:  "Champion",
		ArtistName: "ArtistName",
	}
	ravenbloomConservatory = database.Card{Name: "Ravenbloom Conservatory", Cost: 1, Type: "Landmark", RegionRefs: []string{"Noxus"}}
	theDarkinBallista      = database.Card{Name: "The Darkin Ballista", Cost: 1, Type: "Equipment", RegionRefs: []string{"Noxus"}}

	crimsonPigeon = database.Card{
		Name:           "Crimson Pigeon",
		Cost:           1,
		Type:           "Unit",
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

	bladesEdge = database.Card{
		Name:           "Blade's Edge",
		Cost:           1,
		Type:           "Spell",
		RegionRefs:     []string{"Noxus"},
		RarityRef:      "Common",
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
