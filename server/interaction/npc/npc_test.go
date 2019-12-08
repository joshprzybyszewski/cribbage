package npc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshprzybyszewski/cribbage/model"
)

func createPlayer(t *testing.T, pID model.PlayerID) *npcPlayer {
	npc, err := NewNPCPlayer(pID, func(a model.PlayerAction) error {
		return nil
	})
	require.Nil(t, err)
	p, ok := npc.(*npcPlayer)
	assert.True(t, ok)
	return p
}

func newNPlayerGame(npcID model.PlayerID, n int) model.Game {
	players := make([]model.Player, n)
	for i := 0; i < n-1; i++ {
		id := model.PlayerID(fmt.Sprintf(`p%d`, i))
		players[i] = model.Player{ID: id}
	}
	players[len(players)-1] = model.Player{ID: npcID}

	hands := make(map[model.PlayerID][]model.Card)
	nCards := 6
	switch n {
	case 3, 4:
		nCards = 5
	}
	for _, p := range players {
		hands[p.ID] = make([]model.Card, nCards)
	}
	return model.Game{
		Players: players,
		Hands:   hands,
	}
}

func TestBuildDealAction(t *testing.T) {
	tests := []struct {
		desc string
		npc  model.PlayerID
		g    model.Game
		exp  model.PlayerAction
	}{{
		desc: `test dumb npc`,
		npc:  `dumbNPC`,
	}, {
		desc: `test simple npc`,
		npc:  `simpleNPC`,
	}, {
		desc: `test calculated npc`,
		npc:  `calculatedNPC`,
	}}
	for _, tc := range tests {
		p := createPlayer(t, tc.npc)

		a := p.buildAction(model.DealCards, tc.g)
		assert.Equal(t, a.Overcomes, model.DealCards)

		da, ok := a.Action.(model.DealAction)
		assert.True(t, ok)
		assert.LessOrEqual(t, da.NumShuffles, 10)
		assert.GreaterOrEqual(t, da.NumShuffles, 1)
	}
}
func TestBuildBuildCribAction(t *testing.T) {
	tests := []struct {
		desc      string
		npc       model.PlayerID
		g         model.Game
		expNCards int
	}{{
		desc:      `test dumb npc`,
		npc:       `dumbNPC`,
		g:         newNPlayerGame(`dumbNPC`, 2),
		expNCards: 2,
	}, {
		desc:      `test simple npc`,
		npc:       `simpleNPC`,
		g:         newNPlayerGame(`simpleNPC`, 2),
		expNCards: 2,
	}, {
		desc:      `test calculated npc`,
		npc:       `calculatedNPC`,
		g:         newNPlayerGame(`calculatedNPC`, 2),
		expNCards: 2,
	}, {
		desc:      `test 3 player game`,
		npc:       `dumbNPC`,
		g:         newNPlayerGame(`dumbNPC`, 3),
		expNCards: 1,
	}, {
		desc:      `test 4 player game`,
		npc:       `dumbNPC`,
		g:         newNPlayerGame(`dumbNPC`, 4),
		expNCards: 1,
	}}
	for _, tc := range tests {
		p := createPlayer(t, tc.npc)

		a := p.buildAction(model.CribCard, tc.g)
		assert.Equal(t, a.Overcomes, model.CribCard)

		bca, ok := a.Action.(model.BuildCribAction)
		assert.True(t, ok)
		assert.Len(t, bca.Cards, tc.expNCards)
	}
}
