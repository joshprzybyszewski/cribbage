package interaction

import (
	"math/rand"

	"github.com/joshprzybyszewski/cribbage/game"
	"github.com/joshprzybyszewski/cribbage/logic/scorer"
	"github.com/joshprzybyszewski/cribbage/model"
	// TODO find a way around this import cycle
	"github.com/joshprzybyszewski/cribbage/server"
)

var npc game.Player

// NewNPCPlayer creates a new NPC with specified type
func NewNPCPlayer(npcType NPC) Player {
	return &npcPlayer{
		Type: npcType,
	}
}

func handleNPCBlocker(npcType NPC, b model.Blocker, g model.Game) error {
	id := model.PlayerID(npcIDs[npcType])
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
		a.Action = handleNPCBuildCrib(npcType, g)
	case model.CutCard:
		a.Action = model.CutDeckAction{
			Percentage: rand.Float64(),
		}
	case model.PegCard:
		a.Action = handleNPCPeg(npcType, g)
	case model.CountHand:
		a.Action = model.CountHandAction{
			Pts: scorer.HandPoints(g.CutCard, g.Hands[id]),
		}
	case model.CountCrib:
		a.Action = model.CountCribAction{
			Pts: scorer.CribPoints(g.CutCard, g.Crib),
		}
	}
	return server.HandleAction(a)
}

func updateNPC(npcType NPC, g model.Game) {
	id := model.PlayerID(npcIDs[npcType])
	switch npcType {
	case Dumb:
		npc = game.NewDumbNPC(g.PlayerColors[id])
	case Simple:
		npc = game.NewSimpleNPC(g.PlayerColors[id])
	case Calculated:
		npc = game.NewCalcNPC(g.PlayerColors[id])
	}
}

func handleNPCPeg(npcType NPC, g model.Game) model.PegAction {
	updateNPC(npcType, g)
	c, sayGo, _ := npc.Peg(g.PeggedCards, g.CurrentPeg())
	return model.PegAction{
		Card:  c,
		SayGo: sayGo,
	}
}

func handleNPCBuildCrib(npcType NPC, g model.Game) model.BuildCribAction {
	updateNPC(npcType, g)
	nCards := 2
	switch len(g.Players) {
	case 3, 4:
		nCards = 1
	}
	return model.BuildCribAction{
		Cards: npc.AddToCrib(g.PlayerColors[g.CurrentDealer], nCards),
	}
}
