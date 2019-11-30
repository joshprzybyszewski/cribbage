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

func handleNPCBlocker(npc *dumbNPCPlayer, b model.Blocker, g model.Game, s string) error {
	a := model.PlayerAction{
		GameID:    g.ID,
		ID:        npc.ID(),
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
			Pts: scorer.HandPoints(g.CutCard, g.Hands[npc.ID()]),
		}
	case model.CountCrib:
		a.Action = model.CountCribAction{
			Pts: scorer.CribPoints(g.CutCard, g.Crib),
		}
	}
	return nil
}
