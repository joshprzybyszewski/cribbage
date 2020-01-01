package interaction

import (
	"github.com/joshprzybyszewski/cribbage/logic/strategy"
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/utils/rand"
)

type npc interface {
	getCribAction([]model.Card, bool) (model.BuildCribAction, error)
	getPegAction([]model.Card, []model.PeggedCard, int) model.PegAction
}

type getCribCards func(desired int, hand []model.Card) []model.Card
type getCribCardsWithErr func(desired int, hand []model.Card) ([]model.Card, error)

// TODO put this in a better place?
// TODO make this function less ugly - should all strategies return errors?
func cribActionHelper(hand []model.Card, npc model.PlayerID, isDealer bool) (model.BuildCribAction, error) {
	var cards []model.Card
	n := len(hand) - 4
	switch npc {
	case Simple:
		if isDealer {
			strats := []getCribCards{
				strategy.GiveCribFifteens,
				strategy.GiveCribPairs,
			}
			cards = strats[rand.Int()%2](n, hand)
		} else {
			strats := []getCribCards{
				strategy.AvoidCribFifteens,
				strategy.AvoidCribPairs,
			}
			cards = strats[rand.Int()%2](n, hand)
		}
	case Calc:
		var err error
		if isDealer {
			strats := []getCribCardsWithErr{
				strategy.KeepHandLowestPotential,
				strategy.GiveCribHighestPotential,
			}
			cards, err = strats[rand.Int()%2](n, hand)
		} else {
			strats := []getCribCardsWithErr{
				strategy.KeepHandHighestPotential,
				strategy.GiveCribLowestPotential,
			}
			cards, err = strats[rand.Int()%2](n, hand)
		}
		if err != nil {
			return model.BuildCribAction{}, err
		}
	}
	return model.BuildCribAction{
		Cards: cards,
	}, nil
}
