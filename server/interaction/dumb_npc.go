package interaction

import (
	"github.com/joshprzybyszewski/cribbage/model"
)

var _ npc = (*dumbNPC)(nil)

type dumbNPC struct{}

func (npc *dumbNPC) getBuildCribAction(hand []model.Card, _ bool) (model.BuildCribAction, error) {
	n := len(hand) - 4
	return model.BuildCribAction{
		Cards: hand[0:n],
	}, nil
}

func (npc *dumbNPC) getPegAction(unpegged []model.Card, _ []model.PeggedCard, curPeg int) model.PegAction {
	maxVal := model.MaxPeggingValue - curPeg
	for _, c := range unpegged {
		if c.PegValue() > maxVal {
			continue
		}
		return model.PegAction{
			Card:  c,
			SayGo: false,
		}

	}
	return model.PegAction{
		Card:  model.Card{},
		SayGo: true,
	}
}
