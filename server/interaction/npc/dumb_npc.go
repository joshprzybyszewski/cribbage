package npc

import (
	"github.com/joshprzybyszewski/cribbage/model"
)

const maxPeggingValue = 31

var _ npcLogic = (*dumbNPCLogic)(nil)

type dumbNPCLogic struct{}

func (npc *dumbNPCLogic) addToCrib(hand []model.Card, _ bool) []model.Card {
	n := len(hand) - 4
	return hand[0:n]
}

func (npc *dumbNPCLogic) peg(hand []model.Card, _ []model.PeggedCard, curPeg int) (model.Card, bool) {
	maxVal := maxPeggingValue - curPeg
	for _, c := range hand {
		if c.PegValue() > maxVal {
			continue
		}
		return c, false
	}
	return model.Card{}, true
}
