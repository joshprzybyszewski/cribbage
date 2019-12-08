package interaction

import (
	"github.com/joshprzybyszewski/cribbage/game"
	"github.com/joshprzybyszewski/cribbage/logic/scorer"
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/utils/rand"
)

// NPCType is an enum specifying which type of NPCType
type NPCType int

// Dumb, Simple, and Calculated are supported
const (
	Dumb NPCType = iota
	Simple
	Calculated
)

var _ Player = (*npcPlayer)(nil)

type npcPlayer struct {
	Type NPCType

	id                   model.PlayerID
	handleActionCallback func(a model.PlayerAction) error
}

var npcs = [...]npcPlayer{
	Dumb: {
		Type: Dumb,
		id:   `dumbNPC`,
	},
	Simple: {
		Type: Simple,
		id:   `simpleNPC`,
	},
	Calculated: {
		Type: Calculated,
		id:   `calculatedNPC`,
	},
}

var me game.Player

// NewNPCPlayer creates a new NPC with specified type
func NewNPCPlayer(n NPCType, cb func(a model.PlayerAction) error) Player {
	npc := npcs[n]
	npc.handleActionCallback = cb
	return &npc
}

func (npc *npcPlayer) ID() model.PlayerID {
	return npc.id
}

func (npc *npcPlayer) NotifyBlocking(b model.Blocker, g model.Game, s string) error {
	a := npc.buildAction(b, g)
	return npc.handleActionCallback(a)
}

// The NPC doesn't care about messages or score updates
func (npc *npcPlayer) NotifyMessage(g model.Game, s string) error {
	return nil
}
func (npc *npcPlayer) NotifyScoreUpdate(g model.Game, msgs ...string) error {
	return nil
}

func (npc *npcPlayer) buildAction(b model.Blocker, g model.Game) model.PlayerAction {
	a := model.PlayerAction{
		GameID:    g.ID,
		ID:        npc.ID(),
		Overcomes: b,
	}
	switch b {
	case model.DealCards:
		a.Action = model.DealAction{
			NumShuffles: rand.Intn(10) + 1,
		}
	case model.CribCard:
		a.Action = npc.handleBuildCrib(g)
	case model.CutCard:
		a.Action = model.CutDeckAction{
			Percentage: rand.Float64(),
		}
	case model.PegCard:
		a.Action = npc.handlePeg(g)
	case model.CountHand:
		a.Action = model.CountHandAction{
			Pts: scorer.HandPoints(g.CutCard, g.Hands[npc.ID()]),
		}
	case model.CountCrib:
		a.Action = model.CountCribAction{
			Pts: scorer.CribPoints(g.CutCard, g.Crib),
		}
	}
	return a
}

func (npc *npcPlayer) updateCurrentNPC(g model.Game) {
	id := npc.ID()
	switch npc.Type {
	case Dumb:
		me = game.NewDumbNPC(g.PlayerColors[id])
	case Simple:
		me = game.NewSimpleNPC(g.PlayerColors[id])
	case Calculated:
		me = game.NewCalcNPC(g.PlayerColors[id])
	}
}

func (npc *npcPlayer) handlePeg(g model.Game) model.PegAction {
	npc.updateCurrentNPC(g)
	c, sayGo, _ := me.Peg(g.PeggedCards, g.CurrentPeg())
	return model.PegAction{
		Card:  c,
		SayGo: sayGo,
	}
}

func (npc *npcPlayer) handleBuildCrib(g model.Game) model.BuildCribAction {
	npc.updateCurrentNPC(g)
	nCards := 2
	switch len(g.Players) {
	case 3, 4:
		nCards = 1
	}
	return model.BuildCribAction{
		Cards: me.AddToCrib(g.PlayerColors[g.CurrentDealer], nCards),
	}
}
