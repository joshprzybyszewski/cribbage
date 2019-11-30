package interaction

import (
	"math/rand"

	"github.com/joshprzybyszewski/cribbage/model"
)

type npcPlayer interface {
	id() model.PlayerID
	buildCrib(g model.Game) model.BuildCribAction
	peg(g model.Game) model.PegAction
	countHand(g model.Game) model.CountHandAction
	countCrib(g model.Game) model.CountCribAction
}

func handleNPCBlocker(npc *dumbNPCPlayer, b model.Blocker, g model.Game, s string) error {
	a := model.PlayerAction{
		GameID:    g.ID,
		ID:        npc.id(),
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
		a.Action = npc.countHand(g)
	case model.CountCrib:
		a.Action = npc.countCrib(g)
	}
	return nil
}
