package strategy

import (
	"github.com/joshprzybyszewski/cribbage/cards"
	"github.com/joshprzybyszewski/cribbage/pegging"
)

func PegHighestCardNow(hand, prevPegs []cards.Card, curPeg int) (cards.Card, bool) {
	bestCard := cards.Card{}
	bestPoints := 0

	for _, c := range hand {
		if curPeg+c.PegValue() > 31 {
			continue
		}

		p, err := pegging.PointsForCard(prevPegs, c)
		if err != nil {
			return cards.Card{}, false
		}

		if p > bestPoints {
			bestCard = c
			bestPoints = p
		}
	}

	return bestCard, bestCard == cards.Card{}
}
