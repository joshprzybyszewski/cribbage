package interaction

import (
	"github.com/joshprzybyszewski/cribbage/game"
	"github.com/joshprzybyszewski/cribbage/logic/scorer"
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/utils/rand"

	// TODO find a way around this import cycle
	"github.com/joshprzybyszewski/cribbage/server"
)

// NPC is an enum specifying which type of NPC
type NPC int

// Dumb, Simple, and Calculated are supported
const (
	Dumb NPC = iota
	Simple
	Calculated
)

var npcIDs = [...]string{
	Dumb:       `dumbNPC`,
	Simple:     `simpleNPC`,
	Calculated: `calculatedNPC`,
}

var _ Player = (*npcPlayer)(nil)

type npcPlayer struct {
	Type NPC
}

func (npc *npcPlayer) ID() model.PlayerID {
	return model.PlayerID(npcIDs[npc.Type])
}

func (npc *npcPlayer) NotifyBlocking(b model.Blocker, g model.Game, s string) error {
	return handleNPCBlocker(npc.Type, b, g)
}
func (npc *npcPlayer) NotifyMessage(g model.Game, s string) error {
	return nil
}
func (npc *npcPlayer) NotifyScoreUpdate(g model.Game, msgs ...string) error {
	return nil
}

// NewNPCPlayer creates a new NPC with specified type
func NewNPCPlayer(npcType NPC) Player {
	return &npcPlayer{
		Type: npcType,
	}
}

var npc game.Player

func handleNPCBlocker(npcType NPC, b model.Blocker, g model.Game) error {
	id := model.PlayerID(npcIDs[npcType])
	action := model.PlayerAction{
		GameID:    g.ID,
		ID:        id,
		Overcomes: b,
	}
	switch b {
	case model.DealCards:
		action.Action = model.DealAction{
			NumShuffles: rand.Intn(10),
		}
	case model.CribCard:
		action.Action = handleNPCBuildCrib(npcType, g)
	case model.CutCard:
		action.Action = model.CutDeckAction{
			Percentage: rand.Float64(),
		}
	case model.PegCard:
		action.Action = handleNPCPeg(npcType, g)
	case model.CountHand:
		action.Action = model.CountHandAction{
			Pts: scorer.HandPoints(g.CutCard, g.Hands[id]),
		}
	case model.CountCrib:
		action.Action = model.CountCribAction{
			Pts: scorer.CribPoints(g.CutCard, g.Crib),
		}
	}
	return server.HandleAction(action)
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
