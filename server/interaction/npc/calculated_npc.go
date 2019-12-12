package npc

import (
	"github.com/joshprzybyszewski/cribbage/logic/strategy"
	"github.com/joshprzybyszewski/cribbage/model"
)

var _ npcLogic = (*calcNPCLogic)(nil)

type calcNPCLogic struct{}

func (npc *calcNPCLogic) getCribAction(hand []model.Card, isDealer bool) model.BuildCribAction {
	return cribActionHelper(hand, Calc, isDealer)
}

func (npc *calcNPCLogic) getPegAction(hand []model.Card, prevPegs []model.PeggedCard, curPeg int) model.PegAction {
	card, sayGo := strategy.PegHighestCardNow(hand, prevPegs, curPeg)
	return model.PegAction{
		Card:  card,
		SayGo: sayGo,
	}
}
