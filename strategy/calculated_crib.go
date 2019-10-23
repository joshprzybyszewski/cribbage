package strategy

import (
	"github.com/joshprzybyszewski/cribbage/cards"
	"github.com/joshprzybyszewski/cribbage/scorer"
)

// GiveCribHighestPotential gives the crib the highest potential pointed crib
func GiveCribHighestPotential(_ int, hand []cards.Card) []cards.Card {
	isBetter := func(old, new float64) bool { return new > old }
	return getBestPotential(hand, isBetter)
}

// GiveCribLowestPotential gives the crib the lowest potential pointed hand
func GiveCribLowestPotential(_ int, hand []cards.Card) []cards.Card {
	isBetter := func(old, new float64) bool { return new < old }
	return getBestPotential(hand, isBetter)
}

func getBestPotential(hand []cards.Card, isBetter func(old, new float64) bool) []cards.Card {
	bestCrib := make([]cards.Card, 0, len(hand)-4)
	bestPotential := 0.0

	allDeposits := getDeposits(hand)

	for i, dep := range allDeposits {
		p := getPotentialForDeposit(dep)
		if i == 0 || isBetter(bestPotential, p) {
			bestCrib = bestCrib[:0]
			bestCrib = append(bestCrib, dep...)
			bestPotential = p
		}
	}

	return bestCrib
}

func getDeposits(hand []cards.Card) [][]cards.Card {
	if len(hand) > 6 || len(hand) <= 4 {
		return nil
	}

	// 6 choose 2 = 15, our largest number of potential deposits to the crib
	allDeposits := make([][]cards.Card, 0, 15)
	for i, c1 := range hand {
		if len(hand) == 5 {
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

func getPotentialForDeposit(myCrib []cards.Card) float64 {
	myCards := map[string]struct{}{}
	for _, c := range myCrib {
		myCards[c.String()] = struct{}{}
	}

	totalCribPoints := 0
	totalHands := 0

	for i := 0; i < 52; i++ {
		lead := cards.NewCardFromNumber(i)
		if _, ok := myCards[lead.String()]; ok {
			continue
		}

		for o1 := 0; o1 < 52; o1++ {
			if i == o1 {
				continue
			}
			oc1 := cards.NewCardFromNumber(o1)
			if _, ok := myCards[oc1.String()]; ok {
				continue
			}
			for o2 := o1 + 1; o2 < 52; o2++ {
				if i == o2 {
					continue
				}
				oc2 := cards.NewCardFromNumber(o2)
				if _, ok := myCards[oc2.String()]; ok {
					continue
				}

				if len(myCards) == 2 {
					totalCribPoints += scorer.CribPoints(lead, append(myCrib, oc1, oc2))
					totalHands++
					continue
				}

				for o3 := o2 + 1; o3 < 52; o3++ {
					if i == o3 {
						continue
					}

					oc3 := cards.NewCardFromNumber(o3)
					if _, ok := myCards[oc3.String()]; ok {
						continue
					}

					totalCribPoints += scorer.CribPoints(lead, append(myCrib, oc1, oc2, oc3))
					totalHands++
				}
			}
		}
	}

	return float64(totalCribPoints) / float64(totalHands)
}
