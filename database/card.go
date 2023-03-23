package database

// Card represents the card structure present in LoR Data Dragon
type Card struct {
	AssociatedCards    []string `json:"associatedCards"`
	AssociatedCardRefs []string `json:"associatedCardRefs"`
	Assets             []struct {
		GameAbsolutePath string `json:"gameAbsolutePath"`
		FullAbsolutePath string `json:"fullAbsolutePath"`
	} `json:"assets"`
	Regions               []string `json:"regions"`
	RegionRefs            []string `json:"regionRefs"`
	Attack                int      `json:"attack"`
	Cost                  int      `json:"cost"`
	Health                int      `json:"health"`
	Description           string   `json:"description"`
	DescriptionRaw        string   `json:"descriptionRaw"`
	LevelupDescription    string   `json:"levelupDescription"`
	LevelupDescriptionRaw string   `json:"levelupDescriptionRaw"`
	FlavorText            string   `json:"flavorText"`
	ArtistName            string   `json:"artistName"`
	Name                  string   `json:"name"`
	CardCode              string   `json:"cardCode"`
	Keywords              []string `json:"keywords"`
	KeywordRefs           []string `json:"keywordRefs"`
	SpellSpeed            string   `json:"spellSpeed"`
	SpellSpeedRef         string   `json:"spellSpeedRef"`
	Rarity                string   `json:"rarity"`
	RarityRef             string   `json:"rarityRef"`
	Subtypes              []string `json:"subtypes"`
	Supertype             string   `json:"supertype"`
	Type                  string   `json:"type"`
	Collectible           bool     `json:"collectible"`
	Set                   string   `json:"set"`
}

type cards []Card

func (cs cards) String(i int) string {
	return cs[i].Name
}

func (cs cards) Len() int {
	return len(cs)
}
