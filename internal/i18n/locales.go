package i18n

import (
	"embed"
	"encoding/json"
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed *.json
var static embed.FS

type Locale string

var Default Locale = "en_us"

func (l Locale) Name() string {
	return map[Locale]string{
		"de_de": "Deutsch",
		"en_us": "English",
		"es_es": "Español (EUW)",
		"es_mx": "Español (LATAM)",
		"fr_fr": "Français",
		"it_it": "Italiano",
		"ja_jp": "日本語",
		"ko_kr": "한국어",
		"pl_pl": "Polski",
		"pt_br": "Português",
		"ru_ru": "Русский",
		"th_th": "ภาษาไทย",
		"tr_tr": "Türkçe",
		"vi_vn": "Tiếng Việt",
		"zh_tw": "繁體中文",
	}[l]

}

var Locales = []Locale{
	"de_de",
	"en_us",
	"es_es",
	"es_mx",
	"fr_fr",
	"it_it",
	"ja_jp",
	"ko_kr",
	"pl_pl",
	"pt_br",
	"ru_ru",
	"th_th",
	"tr_tr",
	"vi_vn",
	"zh_tw",
}

func AsStringSlice(ls []Locale) []string {
	strs := make([]string, len(ls))
	for i, l := range ls {
		strs[i] = string(l)
	}

	return strs
}

var DefaultLoc = LoadTranslations()

type Localizer struct {
	bundle *i18n.Bundle
}

func (l Localizer) Localize(language string, messageID string) string {
	localizer := i18n.NewLocalizer(l.bundle, language)
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
}

func LoadTranslations() *Localizer {
	t, err := language.Parse("en_us")
	if err != nil {
		fmt.Println(err)
	}
	bundle := i18n.NewBundle(t)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	for _, l := range Locales {
		_, err := bundle.LoadMessageFileFS(static, fmt.Sprintf("messages.%s.json", l))
		if err != nil {
			fmt.Println(err)
		}
	}

	return &Localizer{bundle: bundle}
}
