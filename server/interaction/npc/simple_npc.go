package npc

import (
	"github.com/joshprzybyszewski/cribbage/logic/strategy"
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/utils/rand"
)

var _ npcLogic = (*simpleNPCLogic)(nil)

type simpleNPCLogic struct{}

func (npc *simpleNPCLogic) getCribAction(hand []model.Card, isDealer bool) model.BuildCribAction {
	dealerStrats := []func(desired int, hand []model.Card) []model.Card{
		strategy.GiveCribFifteens,
		strategy.GiveCribPairs,
	}
	notDealerStrats := []func(desired int, hand []model.Card) []model.Card{
		strategy.AvoidCribFifteens,
		strategy.AvoidCribPairs,
	}
	return cribActionHelper(hand, isDealer, dealerStrats, notDealerStrats)
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
