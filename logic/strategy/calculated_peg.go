package strategy

import (
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/logic/pegging"
)

func PegHighestCardNow(hand, prevPegs []model.Card, curPeg int) (model.Card, bool) {
	bestCard := model.Card{}
	bestPoints := 0

	for _, c := range hand {
		if curPeg+c.PegValue() > 31 {
			continue
		}

		p, err := pegging.PointsForCard(prevPegs, c)
		if err != nil {
			return model.Card{}, false
		}

		if p > bestPoints {
			bestCard = c
			bestPoints = p
		}
	}

	return bestCard, bestCard == model.Card{}
}
