package interaction

import (
	"math/rand"

	"github.com/joshprzybyszewski/cribbage/logic/scorer"
	"github.com/joshprzybyszewski/cribbage/model"
)

type npcPlayer interface {
	buildCrib(g model.Game) model.BuildCribAction
	peg(g model.Game) model.PegAction
}

func handleNPCBlocker(npc npcPlayer, b model.Blocker, g model.Game, id model.PlayerID) error {
	a := model.PlayerAction{
		GameID:    g.ID,
		ID:        id,
		Overcomes: b,
	}
	switch b {
	case model.DealCards:
		a.Action = model.DealAction{
			NumShuffles: rand.Intn(10),
		}
	case model.CribCard:
		a.Action = npc.buildCrib(g)
	case model.CutCard:
		a.Action = model.CutDeckAction{
			Percentage: rand.Float64(),
		}
	case model.PegCard:
		a.Action = npc.peg(g)
	case model.CountHand:
		a.Action = model.CountHandAction{
			Pts: scorer.HandPoints(g.CutCard, g.Hands[id]),
		}
	case model.CountCrib:
		a.Action = model.CountCribAction{
			Pts: scorer.CribPoints(g.CutCard, g.Crib),
		}
	}
	return nil
}
