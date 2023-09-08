package card_test

import (
	"testing"

	"github.com/dneto/sai-scout/internal/card"
	"github.com/dneto/sai-scout/internal/repository"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCard(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Decks Suite")
}

var _ = Describe("Card", func() {
	DescribeTable("IsChampion",
		func(c *repository.Card, expected bool) {
			Expect(card.IsChampion(c)).To(Equal(expected))
		},
		Entry("Annie is a champion", annie, true),
		Entry("Ravenbloom Conservatory is not a champion", ravenbloomConservatory, false),
		Entry("The Darkin Balist is not a champion", theDarkinBallista, false),
		Entry("Crimson Pigeon is not a champion", crimsonPigeon, false),
		Entry("Blade's Edge is not a champion", bladesEdge, false),
	)

	DescribeTable("IsFollower",
		func(c *repository.Card, expected bool) {
			Expect(card.IsFollower(c)).To(Equal(expected))
		},
		Entry("Annie is not a follower", annie, false),
		Entry("Ravenbloom Conservatory is not a follower", ravenbloomConservatory, false),
		Entry("The Darkin Balist is not a follower", theDarkinBallista, false),
		Entry("Crimson Pigeon is a follower", crimsonPigeon, true),
		Entry("Blade's Edge is not a follower", bladesEdge, false),
	)

	DescribeTable("IsSpell",
		func(c *repository.Card, expected bool) {
			Expect(card.IsSpell(c)).To(Equal(expected))
		},
		Entry("Annie is not a spell", annie, false),
		Entry("Ravenbloom Conservatory is not a spell", ravenbloomConservatory, false),
		Entry("The Darkin Balist is not a spell", theDarkinBallista, false),
		Entry("Crimson Pigeon is not a spell", crimsonPigeon, false),
		Entry("Blade's Edge is a spell", bladesEdge, true),
	)

	DescribeTable("IsLandmark",
		func(c *repository.Card, expected bool) {
			Expect(card.IsLandmark(c)).To(Equal(expected))
		},
		Entry("Annie is not a landmark", annie, false),
		Entry("Ravenbloom Conservatory is a landmark", ravenbloomConservatory, true),
		Entry("The Darkin Balist is not a landmark", theDarkinBallista, false),
		Entry("Crimson Pigeon is not a landmark", crimsonPigeon, false),
		Entry("Blade's Edge is not a spell landmark", bladesEdge, false),
	)

	DescribeTable("IsEquipment",
		func(c *repository.Card, expected bool) {
			Expect(card.IsLandmark(c)).To(Equal(expected))
		},
		Entry("Annie is not an equipment", annie, false),
		Entry("Ravenbloom Conservatory is not an equipment", ravenbloomConservatory, true),
		Entry("The Darkin Balist is an equipment", theDarkinBallista, false),
		Entry("Crimson Pigeon is not an equipment", crimsonPigeon, false),
		Entry("Blade's Edge is not an equipment", bladesEdge, false),
	)
})

var (
	annie                  = &repository.Card{CardCode: "06NX012", Name: "Annie", Cost: 0, RarityRef: "Champion", TypeRef: "Unit"}
	ravenbloomConservatory = &repository.Card{CardCode: "06NX028", Name: "Ravenbloom Conservatory", Cost: 1, TypeRef: "Landmark"}
	theDarkinBallista      = &repository.Card{CardCode: "06NX020", Name: "The Darkin Ballista", Cost: 2, TypeRef: "Equipment"}
	crimsonPigeon          = &repository.Card{CardCode: "06NX041", Name: "Crimson Pigeon", Cost: 3, TypeRef: "Unit"}
	bladesEdge             = &repository.Card{CardCode: "01NX043", Name: "Blade's Edge", Cost: 4, TypeRef: "Spell"}
)
