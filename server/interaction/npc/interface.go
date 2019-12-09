package npc

import (
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/utils/rand"
)

type npcLogic interface {
	getCribAction([]model.Card, bool) model.BuildCribAction
	getPegAction([]model.Card, []model.PeggedCard, int) model.PegAction
}

// TODO put this in a better place?
func cribActionHelper(hand []model.Card,
	isDealer bool, dealerStrats,
	notDealerStrats []func(int, []model.Card) []model.Card) model.BuildCribAction {

	var i int
	n := len(hand) - 4
	if isDealer {
		i = rand.Int() % len(dealerStrats)
		return model.BuildCribAction{
			Cards: dealerStrats[i](n, hand),
		}
	}
	i = rand.Int() % len(notDealerStrats)
	return model.BuildCribAction{
		Cards: notDealerStrats[i](n, hand),
	}
}
