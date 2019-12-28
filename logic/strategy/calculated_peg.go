package strategy

import (
	"github.com/joshprzybyszewski/cribbage/logic/pegging"
	"github.com/joshprzybyszewski/cribbage/model"
)

func PegHighestCardNow(hand []model.Card, prevPegs []model.PeggedCard, curPeg int) (model.Card, bool) {
	bestCard := model.Card{}
	bestPoints := -1
	cardsOverMax := 0

	for _, c := range hand {
		if curPeg+c.PegValue() > 31 {
			cardsOverMax++
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
	if cardsOverMax == len(hand) {
		return model.Card{}, true
	}
	return bestCard, false
}
