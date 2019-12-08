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

func newNPlayerGame(npcID model.PlayerID, n int, pegCards []model.Card) model.Game {
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
	hands[npcID] = []model.Card{
		model.NewCardFromString(`6c`),
		model.NewCardFromString(`7c`),
		model.NewCardFromString(`8c`),
		model.NewCardFromString(`9c`),
	}

	pegs := make([]model.PeggedCard, 0)
	for i, c := range pegCards {
		pegs = append(pegs, model.PeggedCard{
			Card:     c,
			PlayerID: players[i%n].ID,
		})
	}
	return model.Game{
		Players:     players,
		Hands:       hands,
		PeggedCards: pegs,
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
func TestBuildCutAction(t *testing.T) {
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

		a := p.buildAction(model.CutCard, tc.g)
		assert.Equal(t, a.Overcomes, model.CutCard)

		cda, ok := a.Action.(model.CutDeckAction)
		assert.True(t, ok)
		assert.LessOrEqual(t, cda.Percentage, 1.0)
		assert.GreaterOrEqual(t, cda.Percentage, 0.0)
	}
}
func TestPegAction(t *testing.T) {
	tests := []struct {
		desc  string
		npc   model.PlayerID
		g     model.Game
		expGo bool
	}{{
		desc:  `test dumb npc`,
		npc:   `dumbNPC`,
		g:     newNPlayerGame(`dumbNPC`, 2, make([]model.Card, 0)),
		expGo: false,
	}, {
		desc:  `test simple npc`,
		npc:   `simpleNPC`,
		g:     newNPlayerGame(`simpleNPC`, 2, make([]model.Card, 0)),
		expGo: false,
	}, {
		desc:  `test calculated npc`,
		npc:   `calculatedNPC`,
		g:     newNPlayerGame(`calculatedNPC`, 2, make([]model.Card, 0)),
		expGo: false,
	}, {
		desc: `test go`,
		npc:  `dumbNPC`,
		g: newNPlayerGame(`dumbNPC`, 2, []model.Card{
			model.NewCardFromString(`10c`),
			model.NewCardFromString(`10s`),
			model.NewCardFromString(`10h`),
		}),
		expGo: true,
	}}
	for _, tc := range tests {
		p := createPlayer(t, tc.npc)

		a := p.buildAction(model.PegCard, tc.g)
		assert.Equal(t, a.Overcomes, model.PegCard)

		pa, ok := a.Action.(model.PegAction)
		assert.True(t, ok)
		if tc.expGo {
			assert.True(t, pa.SayGo)
		} else {
			assert.False(t, pa.SayGo, tc.desc)
			assert.NotEqual(t, model.Card{}, pa.Card)
		}
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
		g:         newNPlayerGame(`dumbNPC`, 2, make([]model.Card, 0)),
		expNCards: 2,
	}, {
		desc:      `test simple npc`,
		npc:       `simpleNPC`,
		g:         newNPlayerGame(`simpleNPC`, 2, make([]model.Card, 0)),
		expNCards: 2,
	}, {
		desc:      `test calculated npc`,
		npc:       `calculatedNPC`,
		g:         newNPlayerGame(`calculatedNPC`, 2, make([]model.Card, 0)),
		expNCards: 2,
	}, {
		desc:      `test 3 player game`,
		npc:       `dumbNPC`,
		g:         newNPlayerGame(`dumbNPC`, 3, make([]model.Card, 0)),
		expNCards: 1,
	}, {
		desc:      `test 4 player game`,
		npc:       `dumbNPC`,
		g:         newNPlayerGame(`dumbNPC`, 4, make([]model.Card, 0)),
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
