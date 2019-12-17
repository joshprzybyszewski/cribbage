package npc

import (
	"github.com/joshprzybyszewski/cribbage/logic/strategy"
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/utils/rand"
)

type npcLogic interface {
	getCribAction([]model.Card, bool) (model.BuildCribAction, error)
	getPegAction([]model.Card, []model.PeggedCard, int) model.PegAction
}

type getCribCards func(desired int, hand []model.Card) ([]model.Card, error)

// TODO put this in a better place?
func cribActionHelper(hand []model.Card, npc model.PlayerID, isDealer bool) (model.BuildCribAction, error) {
	var strats []getCribCards
	switch npc {
	case `simpleNPC`:
		if isDealer {
			strats = []getCribCards{
				func(d int, h []model.Card) ([]model.Card, error) {
					return strategy.GiveCribFifteens(d, h), nil
				},
				func(d int, h []model.Card) ([]model.Card, error) {
					return strategy.GiveCribPairs(d, h), nil
				}}
		} else {
			strats = []getCribCards{
				func(d int, h []model.Card) ([]model.Card, error) {
					return strategy.AvoidCribFifteens(d, h), nil
				},
				func(d int, h []model.Card) ([]model.Card, error) {
					return strategy.AvoidCribPairs(d, h), nil
				}}
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
	cards, err := strats[i](n, hand)
	if err != nil {
		return model.BuildCribAction{}, err
	}
	return model.BuildCribAction{
		Cards: cards,
	}, nil
}
