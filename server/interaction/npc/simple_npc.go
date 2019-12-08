package npc

import (
	"github.com/joshprzybyszewski/cribbage/logic/strategy"
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/utils/rand"
)

var _ npcLogic = (*simpleNPCLogic)(nil)

type simpleNPCLogic struct{}

func (npc *simpleNPCLogic) addToCrib(hand []model.Card, isDealer bool) []model.Card {
	n := len(hand) - 4
	if isDealer {
		if rand.Int()%2 == 0 {
			return strategy.GiveCribFifteens(n, hand)
		}
		return strategy.GiveCribPairs(n, hand)
	}

	if rand.Int()%2 == 0 {
		return strategy.AvoidCribFifteens(n, hand)
	}
	return strategy.AvoidCribPairs(n, hand)
}

func (npc *simpleNPCLogic) peg(hand []model.Card, prevPegs []model.PeggedCard, curPeg int) (model.Card, bool) {
	switch rand.Int() % 4 {
	case 0:
		return strategy.PegToFifteen(hand, prevPegs, curPeg)
	case 1:
		return strategy.PegToThirtyOne(hand, prevPegs, curPeg)
	case 2:
		return strategy.PegToPair(hand, prevPegs, curPeg)
	}
	return strategy.PegToRun(hand, prevPegs, curPeg)
}
