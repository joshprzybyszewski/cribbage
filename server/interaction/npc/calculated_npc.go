package npc

import (
	"github.com/joshprzybyszewski/cribbage/logic/strategy"
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/utils/rand"
)

var _ npcLogic = (*calcNPCLogic)(nil)

type calcNPCLogic struct{}

func (npc *calcNPCLogic) addToCrib(g model.Game, pID model.PlayerID, n int) []model.Card {
	hand := g.Hands[pID]
	if pID == g.CurrentDealer {
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
