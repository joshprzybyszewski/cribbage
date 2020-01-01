package interaction

import (
	"context"
	"errors"

	"github.com/joshprzybyszewski/cribbage/logic/scorer"
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/utils/rand"
)

const (
	Dumb   model.PlayerID = `DumbNPC`
	Simple model.PlayerID = `SimpleNPC`
	Calc   model.PlayerID = `CalculatedNPC`
)

var (
	ErrUnknownNPCType = errors.New(`unknown NPC type`)
)

var npcs = map[model.PlayerID]npc{
	Dumb:   &dumbNPC{},
	Simple: &simpleNPC{},
	Calc:   &calculatedNPC{},
}

var _ Player = (*NPCPlayer)(nil)

type NPCPlayer struct {
	HandleActionCallback func(ctx context.Context, a model.PlayerAction) error
	player               npc
	id                   model.PlayerID
}

// NewNPCPlayer creates a new NPC with specified type
func NewNPCPlayer(pID model.PlayerID, cb func(ctx context.Context, a model.PlayerAction) error) (Player, error) {
	p, ok := npcs[pID]
	if !ok {
		return nil, ErrUnknownNPCType
	}
	return &NPCPlayer{
		player:               p,
		id:                   pID,
		HandleActionCallback: cb,
	}, nil
}

func (npc *NPCPlayer) ID() model.PlayerID {
	return npc.id
}

func (npc *NPCPlayer) NotifyBlocking(b model.Blocker, g model.Game, s string) error {
	a, err := npc.buildAction(b, g)
	if err != nil {
		return err
	}
	// TODO find a better way to solve this problem...
	go func() {
		npc.HandleActionCallback(context.Background(), a)
	}()
	return nil
}

// The NPC doesn't care about messages or score updates
func (npc *NPCPlayer) NotifyMessage(g model.Game, s string) error {
	return nil
}
func (npc *NPCPlayer) NotifyScoreUpdate(g model.Game, msgs ...string) error {
	return nil
}

func getUnpeggedCards(hand []model.Card, pc []model.PeggedCard) []model.Card {
	peggedMap := make(map[model.Card]struct{}, len(pc))
	cardsLeft := make([]model.Card, 0, len(hand))
	for _, c := range pc {
		peggedMap[c.Card] = struct{}{}
	}
	for _, c := range hand {
		if _, ok := peggedMap[c]; ok {
			continue
		}
		cardsLeft = append(cardsLeft, c)
	}
	return cardsLeft
}

func (npc *NPCPlayer) buildAction(b model.Blocker, g model.Game) (model.PlayerAction, error) {
	a := model.PlayerAction{
		GameID:    g.ID,
		ID:        npc.ID(),
		Overcomes: b,
	}
	myHand := g.Hands[npc.ID()]
	switch b {
	case model.DealCards:
		a.Action = model.DealAction{
			NumShuffles: rand.Intn(10) + 1,
		}
	case model.CribCard:
		bca, err := npc.player.getBuildCribAction(myHand, g.CurrentDealer == npc.ID())
		if err != nil {
			return model.PlayerAction{}, err
		}
		a.Action = bca
	case model.CutCard:
		a.Action = model.CutDeckAction{
			Percentage: rand.Float64(),
		}
	case model.PegCard:
		cardsLeft := getUnpeggedCards(myHand, g.PeggedCards)
		a.Action = npc.player.getPegAction(cardsLeft, g.PeggedCards, g.CurrentPeg())
	case model.CountHand:
		a.Action = model.CountHandAction{
			Pts: scorer.HandPoints(g.CutCard, myHand),
		}
	case model.CountCrib:
		a.Action = model.CountCribAction{
			Pts: scorer.CribPoints(g.CutCard, g.Crib),
		}
	}
	return a, nil
}
