package npc

import (
	"github.com/joshprzybyszewski/cribbage/logic/strategy"
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/utils/rand"
)

var _ npcLogic = (*simpleNPCInteraction)(nil)

type simpleNPCInteraction struct {
	myColor model.PlayerColor
}

func (npc *simpleNPCInteraction) addToCrib(g model.Game, pID model.PlayerID, n int) []model.Card {
	hand := g.Hands[pID]
	if pID == g.CurrentDealer {
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

func (npc *simpleNPCInteraction) peg(g model.Game, pID model.PlayerID) (model.Card, bool) {
	hand := g.Hands[pID]
	prevPegs := g.PeggedCards
	curPeg := g.CurrentPeg()
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
