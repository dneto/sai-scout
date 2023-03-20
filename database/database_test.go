package database_test

import (
	"github.com/dneto/sai-scout/database"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Decoder", func() {
	Context("NewInMemory", func() {
		var (
			target *database.InMemory
			err    error
		)

		Context("receives json data corresponding to LOR Data Dragon", func() {
			BeforeEach(func() {
				target, err = database.NewInMemory(jsonData)

			})

			It("creates the in-memory database without errors", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(target).ToNot(BeNil())
			})
		})

		Context("receives invalid json", func() {
			BeforeEach(func() {
				target, err = database.NewInMemory([]byte(""))

			})

			It("creates the in-memory database without errors", func() {
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("unexpected end of JSON input"))
				Expect(target).To(BeNil())
			})

		})
	})

	Context("CardByCode", func() {
		var (
			card database.Card
			err  error
		)
		Context("receives a existent card code", func() {

			BeforeEach(func() {
				var target *database.InMemory
				target, err = database.NewInMemory(jsonData)
				card, err = target.CardByCode("01IO012")

			})

			It("returns correct card for the given code", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(card.Name).To(Equal("Twin Disciplines"))
			})
		})

		Context("receive a nonexistent card code", func() {
			BeforeEach(func() {
				var target *database.InMemory
				target, err = database.NewInMemory(jsonData)
				card, err = target.CardByCode("01IO001")

			})

			It("returns 'not found' error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("not found"))
			})

		})
	})
})

var jsonData = []byte(`[{
			"associatedCards": [],
			"associatedCardRefs": [],
			"assets": [
			  {
				"gameAbsolutePath": "http://dd.b.pvp.net/4_2_0/set1/en_us/img/cards/01IO012.png",
				"fullAbsolutePath": "http://dd.b.pvp.net/4_2_0/set1/en_us/img/cards/01IO012-full.png"
			  }
			],
			"regions": [
			  "Ionia"
			],
			"regionRefs": [
			  "Ionia"
			],
			"attack": 0,
			"cost": 2,
			"health": 0,
			"description": "Give an ally +2|+0 or +0|+3 this round.",
			"descriptionRaw": "Give an ally +2|+0 or +0|+3 this round.",
			"levelupDescription": "",
			"levelupDescriptionRaw": "",
			"flavorText": "\"Never fear change. It will question you, test your limits. It is our greatest teacher.\" - Karma",
			"artistName": "SIXMOREVODKA",
			"name": "Twin Disciplines",
			"cardCode": "01IO012",
			"keywords": [
			  "Burst"
			],
			"keywordRefs": [
			  "Burst"
			],
			"spellSpeed": "Burst",
			"spellSpeedRef": "Burst",
			"rarity": "COMMON",
			"rarityRef": "Common",
			"subtypes": [],
			"supertype": "",
			"type": "Spell",
			"collectible": true,
			"set": "Set1"
		  }]`,
)
