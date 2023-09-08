package card

import (
	"github.com/dneto/sai-scout/internal/repository"
)

func IsChampion(c *repository.Card) bool {
	return c.RarityRef == "Champion"
}

func IsFollower(c *repository.Card) bool {
	return c.TypeRef == "Unit" && c.RarityRef != "Champion"
}

func IsSpell(c *repository.Card) bool {
	return c.TypeRef == "Spell"
}

func IsLandmark(c *repository.Card) bool {
	return c.TypeRef == "Landmark"
}

func IsEquipment(c *repository.Card) bool {
	return c.TypeRef == "Equipment"
}
