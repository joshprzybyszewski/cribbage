package npc

import (
	"errors"

	"github.com/joshprzybyszewski/cribbage/logic/scorer"
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/utils/rand"
)

var npcs = map[model.PlayerID]npcLogic{
	`dumbNPC`:       &dumbNPCLogic{},
	`simpleNPC`:     &simpleNPCLogic{},
	`calculatedNPC`: &simpleNPCLogic{},
}

var _ interaction.Player = (*npcPlayer)(nil)

type npcPlayer struct {
	logic                npcLogic
	id                   model.PlayerID
	handleActionCallback func(a model.PlayerAction) error
}

// NewNPCPlayer creates a new NPC with specified type
func NewNPCPlayer(pID model.PlayerID, cb func(a model.PlayerAction) error) (interaction.Player, error) {
	l, ok := npcs[pID]
	if !ok {
		return &npcPlayer{}, errors.New(`not a valid npc mode`)
	}
	return &npcPlayer{
		logic:                l,
		id:                   pID,
		handleActionCallback: cb,
	}, nil
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

func (npc *npcPlayer) handlePeg(g model.Game) model.PegAction {
	c, sayGo := npc.logic.peg(g, npc.ID())
	return model.PegAction{
		Card:  c,
		SayGo: sayGo,
	}
}

func (npc *npcPlayer) handleBuildCrib(g model.Game) model.BuildCribAction {
	nCards := len(g.Hands[npc.ID()]) - 4
	return model.BuildCribAction{
		Cards: npc.logic.addToCrib(g, npc.ID(), nCards),
	}
}
