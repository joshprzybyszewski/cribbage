package npc

import (
	"github.com/joshprzybyszewski/cribbage/logic/strategy"
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/utils/rand"
)

type npcLogic interface {
	getCribAction([]model.Card, bool) model.BuildCribAction
	getPegAction([]model.Card, []model.PeggedCard, int) model.PegAction
}

type getCribCards func(desired int, hand []model.Card) []model.Card

// TODO put this in a better place?
func cribActionHelper(hand []model.Card, npc model.PlayerID, isDealer bool) model.BuildCribAction {
	var strats []getCribCards
	switch npc {
	case `simpleNPC`:
		if isDealer {
			strats = []getCribCards{
				strategy.GiveCribFifteens,
				strategy.GiveCribPairs}
		} else {
			strats = []getCribCards{
				strategy.AvoidCribFifteens,
				strategy.AvoidCribPairs}
		}
	case `calculatedNPC`:
		if isDealer {
			strats = []getCribCards{
				strategy.KeepHandLowestPotential,
				strategy.GiveCribHighestPotential}
		} else {
			strats = []getCribCards{
				strategy.KeepHandHighestPotential,
				strategy.GiveCribLowestPotential}
		}
	}

	n := len(hand) - 4
	i := rand.Int() % len(strats)
	return model.BuildCribAction{
		Cards: strats[i](n, hand),
	}
}
