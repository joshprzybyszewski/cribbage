package npc

import (
	"github.com/joshprzybyszewski/cribbage/logic/strategy"
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/utils/rand"
)

var _ npcLogic = (*calcNPCLogic)(nil)

type calcNPCLogic struct{}

func (npc *calcNPCLogic) addToCrib(hand []model.Card, isDealer bool) []model.Card {
	n := len(hand) - 4
	if isDealer {
		if rand.Int()%2 == 0 {
			return strategy.KeepHandLowestPotential(n, hand)
		}
		return strategy.GiveCribHighestPotential(n, hand)
	}

	if rand.Int()%2 == 0 {
		return strategy.KeepHandHighestPotential(n, hand)
	}
	return strategy.GiveCribLowestPotential(n, hand)
}

func (npc *calcNPCLogic) peg(hand []model.Card, prevPegs []model.PeggedCard, curPeg int) (model.Card, bool) {
	return strategy.PegHighestCardNow(hand, prevPegs, curPeg)
}
