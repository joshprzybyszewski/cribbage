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

func chooseFrom(desired int, hand []cards.Card) [][]cards.Card {
	if desired > 2 || desired <= 0 {
		return nil
	}
	// 6 choose 2 = 15, our largest number of potential deposits to the crib
	allDeposits := make([][]cards.Card, 0, 15)
	for i, c1 := range hand {
		if desired == 1 {
			// it's a three or four player game, return one card
			allDeposits = append(allDeposits, []cards.Card{c1})
			continue
		}
		for j := i + 1; j < len(hand); j++ {
			c2 := hand[j]
			allDeposits = append(allDeposits, []cards.Card{c1, c2})
		}
	}
	return allDeposits
}
