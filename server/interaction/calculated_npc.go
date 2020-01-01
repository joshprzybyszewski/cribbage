package interaction

import (
	"github.com/joshprzybyszewski/cribbage/logic/strategy"
	"github.com/joshprzybyszewski/cribbage/model"
)

var _ npc = (*calculatedNPC)(nil)

type calculatedNPC struct{}

func (npc *calculatedNPC) getBuildCribAction(hand []model.Card, isDealer bool) (model.BuildCribAction, error) {
	return cribActionHelper(hand, Calc, isDealer)
}

func (npc *calculatedNPC) getPegAction(unpegged []model.Card, prevPegs []model.PeggedCard, curPeg int) model.PegAction {
	card, sayGo := strategy.PegHighestCardNow(unpegged, prevPegs, curPeg)
	return model.PegAction{
		Card:  card,
		SayGo: sayGo,
	}
}
