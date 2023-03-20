package regions

type Region int

const (
	Demacia      Region = 0
	Freljord     Region = 1
	Ionia        Region = 2
	Noxus        Region = 3
	PiltoverZaun Region = 4
	ShadowIsles  Region = 5
	Bilgewater   Region = 6
	Shurima      Region = 7
	Targon       Region = 9
	BandleCity   Region = 10
	Runeterra    Region = 12
	XX           Region = 99
)

var short = map[Region]string{
	Demacia:      "DE",
	Freljord:     "FR",
	Ionia:        "IO",
	Noxus:        "NX",
	PiltoverZaun: "PZ",
	ShadowIsles:  "SI",
	Bilgewater:   "BW",
	Shurima:      "SH",
	Targon:       "MT",
	BandleCity:   "BC",
	Runeterra:    "RU",
	XX:           "XX",
}

var regionRef = map[Region]string{
	Demacia:      "Demacia",
	Freljord:     "Freljord",
	Ionia:        "Ionia",
	Noxus:        "Noxus",
	PiltoverZaun: "PiltoverZaun",
	ShadowIsles:  "ShadowIsles",
	Bilgewater:   "Bilgewater",
	Shurima:      "Shurima",
	Targon:       "Targon",
	BandleCity:   "BandleCity",
	Runeterra:    "Runeterra",
	XX:           "XX",
}

var emote = map[Region]string{
	Demacia:      "<:de:1086330563151024239>",
	Freljord:     "<:fr:1086325335508926495>",
	Noxus:        "<:nx:1086330895214071921>",
	Ionia:        "<:io:1086330367771934790>",
	PiltoverZaun: "<:pz:1086330622022258738>",
	ShadowIsles:  "<:si:1086329487337214145>",
	Bilgewater:   "<:bw:1086330365154697278>",
	Shurima:      "<:sh:1086330373333602425>",
	Targon:       "<:mt:1086330369420308511>",
	BandleCity:   "<:bc:1086330363841876058>",
	Runeterra:    "<:ru:1086330371844624525>",
}

func (r *Region) Short() string {
	return short[*r]
}

func FromShort(str string) *Region {
	for r, s := range short {
		if s == str {
			return &r
		}
	}

	xx := XX
	return &xx
}

func (r *Region) String() string {
	return regionRef[*r]
}

func (r *Region) Emote() string {
	return emote[*r]
}

func FromString(str string) *Region {
	for r, s := range regionRef {
		if s == str {
			return &r
		}
	}

	xx := XX
	return &xx
}

func Short(str string) string {
	return FromString(str).Short()
}

func Emote(str string) string {
	return FromString(str).Emote()
}
