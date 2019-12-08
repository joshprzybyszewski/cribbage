package npc

import (
	"github.com/joshprzybyszewski/cribbage/logic/strategy"
	"github.com/joshprzybyszewski/cribbage/model"
)

var _ npcLogic = (*calcNPCLogic)(nil)

type calcNPCLogic struct{}

func (npc *calcNPCLogic) getCribAction(hand []model.Card, isDealer bool) model.BuildCribAction {
	dealerStrats := []func(desired int, hand []model.Card) []model.Card{
		strategy.KeepHandLowestPotential,
		strategy.GiveCribHighestPotential,
	}
	notDealerStrats := []func(desired int, hand []model.Card) []model.Card{
		strategy.KeepHandHighestPotential,
		strategy.GiveCribLowestPotential,
	}
	return cribActionHelper(hand, isDealer, dealerStrats, notDealerStrats)
}

func (npc *calcNPCLogic) peg(hand []model.Card, prevPegs []model.PeggedCard, curPeg int) (model.Card, bool) {
	return strategy.PegHighestCardNow(hand, prevPegs, curPeg)
}
