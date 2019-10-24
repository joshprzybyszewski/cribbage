package strategy

import (
	"github.com/joshprzybyszewski/cribbage/cards"
)

func otherOptions(desired int, avoid map[cards.Card]struct{}) [][]cards.Card {
	options := [][]cards.Card{}

	for o1 := 0; o1 < 52; o1++ {
		oc1 := cards.NewCardFromNumber(o1)
		if _, ok := avoid[oc1]; ok {
			continue
		}
		for o2 := o1 + 1; o2 < 52; o2++ {
			oc2 := cards.NewCardFromNumber(o2)
			if _, ok := avoid[oc2]; ok {
				continue
			}

			if desired == 2 {
				options = append(options, []cards.Card{oc1, oc2})
				continue
			}

			for o3 := o2 + 1; o3 < 52; o3++ {
				oc3 := cards.NewCardFromNumber(o3)
				if _, ok := avoid[oc3]; ok {
					continue
				}
				options = append(options, []cards.Card{oc1, oc2, oc3})
			}
		}
	}

	return options
}
