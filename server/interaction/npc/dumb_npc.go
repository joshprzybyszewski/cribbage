package npc

import (
	"github.com/joshprzybyszewski/cribbage/model"
)

const maxPeggingValue = 31

var _ npcLogic = (*dumbNPCLogic)(nil)

type dumbNPCLogic struct{}

func (npc *dumbNPCLogic) addToCrib(g model.Game, pID model.PlayerID, n int) []model.Card {
	return g.Hands[pID][0:n]
}

func (npc *dumbNPCLogic) peg(g model.Game, pID model.PlayerID) (model.Card, bool) {
	hand := g.Hands[pID]
	maxVal := maxPeggingValue - g.CurrentPeg()
	for _, c := range hand {
		if c.PegValue() > maxVal {
			continue
		}
		return c, false
	}
	return model.Card{}, true
}
