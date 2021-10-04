package interaction

import (
	"errors"
	"fmt"
	"time"

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
	actionHandler ActionHandler
	id            model.PlayerID
	player        npc
}

// NewNPCPlayer creates a new NPC with specified type
func NewNPCPlayer(pID model.PlayerID, ah ActionHandler) (Player, error) {
	p, ok := npcs[pID]
	if !ok {
		return nil, ErrUnknownNPCType
	}
	return &NPCPlayer{
		player:        p,
		id:            pID,
		actionHandler: ah,
	}, nil
}

func (npc *NPCPlayer) ID() model.PlayerID {
	return npc.id
}

func (npc *NPCPlayer) NotifyBlocking(b model.Blocker, g model.Game, s string) error {
	pa, err := npc.buildAction(b, g)
	if err != nil {
		return err
	}
	go func() {
		// This is an arbitrary amount of time to sleep. We just need to give
		// the server a chance to increment the phase and get ready to handle
		// our action
		time.Sleep(time.Millisecond * 20)
		err := npc.actionHandler.Handle(pa)
		// TODO do something better with the error...
		if err != nil {
			fmt.Printf("ope! %v\n", err)
		}
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
	pa := model.PlayerAction{
		GameID:       g.ID,
		ID:           npc.ID(),
		Overcomes:    b,
		TimestampStr: time.Now().Format(time.RFC3339),
	}
	myHand := g.Hands[npc.ID()]
	switch b {
	case model.DealCards:
		pa.Action = model.DealAction{
			NumShuffles: rand.Intn(10) + 1,
		}
	case model.CribCard:
		bca, err := npc.player.getBuildCribAction(myHand, g.CurrentDealer == npc.ID())
		if err != nil {
			return model.PlayerAction{}, err
		}
		pa.Action = bca
	case model.CutCard:
		pa.Action = model.CutDeckAction{
			Percentage: rand.Float64(),
		}
	case model.PegCard:
		cardsLeft := getUnpeggedCards(myHand, g.PeggedCards)
		pa.Action = npc.player.getPegAction(cardsLeft, g.PeggedCards, g.CurrentPeg())
	case model.CountHand:
		pa.Action = model.CountHandAction{
			Pts: scorer.HandPoints(g.CutCard, myHand),
		}
	case model.CountCrib:
		pa.Action = model.CountCribAction{
			Pts: scorer.CribPoints(g.CutCard, g.Crib),
		}
	}
	return pa, nil
}
