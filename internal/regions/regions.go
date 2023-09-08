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
	Demacia:      "<:de:1088275845698818109>",
	Freljord:     "<:fr:1088274715405209600>",
	Noxus:        "<:nx:1088276202944479342>",
	Ionia:        "<:io:1088274155755028520>",
	PiltoverZaun: "<:pz:1088276832383668254>",
	ShadowIsles:  "<:si:1088278489356042332>",
	Bilgewater:   "<:bw:1088278952830832671>",
	Shurima:      "<:sh:1088282276753846272>",
	Targon:       "<:mt:1088281324177076244>",
	BandleCity:   "<:bc:1088281710996770878>",
	Runeterra:    "<:ru:1088281714721300511>",
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
