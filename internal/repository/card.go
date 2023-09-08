package repository

type Card struct {
	CardCode string `json:"cardCode"`

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
	Formats               []string `json:"formats"`
	FormatRefs            []string `json:"formatRefs"`

	TypeRef      string
	SupertypeRef string
}
