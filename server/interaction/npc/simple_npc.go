package npc

import (
	"github.com/joshprzybyszewski/cribbage/logic/strategy"
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/utils/rand"
)

var _ npcLogic = (*simpleNPCLogic)(nil)

type simpleNPCLogic struct{}

func (npc *simpleNPCLogic) getCribAction(hand []model.Card, isDealer bool) model.BuildCribAction {
	return cribActionHelper(hand, Simple, isDealer)
}

func (npc *simpleNPCLogic) getPegAction(hand []model.Card, prevPegs []model.PeggedCard, curPeg int) model.PegAction {
	var card model.Card
	var sayGo bool
	switch rand.Int() % 4 {
	case 0:
		card, sayGo = strategy.PegToFifteen(hand, prevPegs, curPeg)
	case 1:
		card, sayGo = strategy.PegToThirtyOne(hand, prevPegs, curPeg)
	case 2:
		card, sayGo = strategy.PegToPair(hand, prevPegs, curPeg)
	default:
		card, sayGo = strategy.PegToRun(hand, prevPegs, curPeg)
	}
	return model.PegAction{
		Card:  card,
		SayGo: sayGo,
	}
}
