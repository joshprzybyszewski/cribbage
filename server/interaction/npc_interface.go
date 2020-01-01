package interaction

import (
	"github.com/joshprzybyszewski/cribbage/logic/strategy"
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/utils/rand"
)

type npc interface {
	getBuildCribAction(hand []model.Card, isDealer bool) (model.BuildCribAction, error)
	getPegAction(unpegged []model.Card, prevPegs []model.PeggedCard, curPeg int) model.PegAction
}

type getCribCards func(desired int, hand []model.Card) ([]model.Card, error)

func cribActionHelper(hand []model.Card, npc model.PlayerID, isDealer bool) (model.BuildCribAction, error) {
	var cards []model.Card
	n := len(hand) - 4
	stratMap := map[model.PlayerID]map[bool][]getCribCards{
		Simple: {
			false: []getCribCards{
				strategy.AvoidCribFifteens,
				strategy.AvoidCribPairs,
			},
			true: []getCribCards{
				strategy.GiveCribFifteens,
				strategy.GiveCribPairs,
			},
		},
		Calc: {
			false: []getCribCards{
				strategy.KeepHandHighestPotential,
				strategy.GiveCribLowestPotential,
			},
			true: []getCribCards{
				strategy.KeepHandLowestPotential,
				strategy.GiveCribHighestPotential,
			},
		},
	}
	strats := stratMap[npc][isDealer]
	idx := rand.Int() % len(strats)
	cards, err := strats[idx](n, hand)
	if err != nil {
		return model.BuildCribAction{}, err
	}
	return model.BuildCribAction{
		Cards: cards,
	}, nil
}
