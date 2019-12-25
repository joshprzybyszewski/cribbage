package npc

import (
	"context"
	"errors"

	"github.com/joshprzybyszewski/cribbage/logic/scorer"
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
	"github.com/joshprzybyszewski/cribbage/utils/rand"
)

const (
	Dumb   model.PlayerID = `dumbNPC`
	Simple model.PlayerID = `simpleNPC`
	Calc   model.PlayerID = `calculatedNPC`
)

var npcs = map[model.PlayerID]npcLogic{
	Dumb:   &dumbNPCLogic{},
	Simple: &simpleNPCLogic{},
	Calc:   &calcNPCLogic{},
}

var _ interaction.Player = (*npcPlayer)(nil)

type npcPlayer struct {
	logic                npcLogic
	id                   model.PlayerID
	handleActionCallback func(ctx context.Context, a model.PlayerAction) error
}

// NewNPCPlayer creates a new NPC with specified type
func NewNPCPlayer(pID model.PlayerID, cb func(ctx context.Context, a model.PlayerAction) error) (interaction.Player, error) {
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
	a, err := npc.buildAction(b, g)
	if err != nil {
		return err
	}
	return npc.handleActionCallback(context.Background(), a)
}

// The NPC doesn't care about messages or score updates
func (npc *npcPlayer) NotifyMessage(g model.Game, s string) error {
	return nil
}
func (npc *npcPlayer) NotifyScoreUpdate(g model.Game, msgs ...string) error {
	return nil
}

func SeedNPCs(db persistence.DB) error {
	npcIDs := []model.PlayerID{Dumb, Simple, Calc}
	for _, id := range npcIDs {
		// we don't need to pass in a callback function when seeding the db
		p, err := NewNPCPlayer(id, nil)
		if err != nil {
			return err
		}
		if _, err := db.GetInteraction(p.ID()); err != nil {
			if err == persistence.ErrPlayerNotFound {
				return db.SaveInteraction(p)
			}
			return err
		}
	}
	return nil
}

func (npc *npcPlayer) buildAction(b model.Blocker, g model.Game) (model.PlayerAction, error) {
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
		act, err := npc.logic.getCribAction(g.Hands[npc.ID()], g.CurrentDealer == npc.ID())
		if err != nil {
			return model.PlayerAction{}, err
		}
		a.Action = act
	case model.CutCard:
		a.Action = model.CutDeckAction{
			Percentage: rand.Float64(),
		}
	case model.PegCard:
		a.Action = npc.logic.getPegAction(g.Hands[npc.ID()], g.PeggedCards, g.CurrentPeg())
	case model.CountHand:
		a.Action = model.CountHandAction{
			Pts: scorer.HandPoints(g.CutCard, g.Hands[npc.ID()]),
		}
	case model.CountCrib:
		a.Action = model.CountCribAction{
			Pts: scorer.CribPoints(g.CutCard, g.Crib),
		}
	}
	return a, nil
}
