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

func newGame(npcID model.PlayerID, n int, pegCards []model.Card) model.Game {
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
	for i := range hands[npcID] {
		// create a hand: 2c, 3c, 4c, ...
		hands[npcID][i] = model.NewCardFromString(fmt.Sprintf(`%dc`, i+2))
	}

	pegs := make([]model.PeggedCard, 0)
	for i, c := range pegCards {
		pegs = append(pegs, model.PeggedCard{
			Card:     c,
			PlayerID: players[i%n].ID,
		})
	}
	return model.Game{
		ID:          5,
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

		a := p.buildAction(model.CutCard, model.Game{})
		assert.Equal(t, a.Overcomes, model.CutCard)

		cda, ok := a.Action.(model.CutDeckAction)
		assert.True(t, ok)
		assert.LessOrEqual(t, cda.Percentage, 1.0)
		assert.GreaterOrEqual(t, cda.Percentage, 0.0)
	}
}
func TestCountHandAction(t *testing.T) {
	g := model.Game{
		CutCard: model.NewCardFromString(`10h`),
	}
	hand := []model.Card{
		model.NewCardFromString(`2c`),
		model.NewCardFromString(`3c`),
		model.NewCardFromString(`4c`),
		model.NewCardFromString(`5c`),
	}
	tests := []struct {
		desc string
		npc  model.PlayerID
		g    model.Game
		exp  model.PlayerAction
	}{{
		desc: `test dumb npc`,
		npc:  Dumb,
		exp: model.PlayerAction{
			ID:        Dumb,
			Overcomes: model.CountHand,
			Action: model.CountHandAction{
				Pts: 12,
			}},
	}, {
		desc: `test simple npc`,
		npc:  Simple,
		exp: model.PlayerAction{
			ID:        Simple,
			Overcomes: model.CountHand,
			Action: model.CountHandAction{
				Pts: 12,
			}},
	}, {
		desc: `test calculated npc`,
		npc:  Calc,
		exp: model.PlayerAction{
			ID:        Calc,
			Overcomes: model.CountHand,
			Action: model.CountHandAction{
				Pts: 12,
			}},
	}}
	for _, tc := range tests {
		p := createPlayer(t, tc.npc)
		g.Hands = map[model.PlayerID][]model.Card{
			tc.npc: hand,
		}

		a := p.buildAction(model.CountHand, g)
		assert.Equal(t, a.Overcomes, tc.exp.Overcomes)

		cha, ok := a.Action.(model.CountHandAction)
		assert.True(t, ok)
		exp, ok := tc.exp.Action.(model.CountHandAction)
		assert.True(t, ok)
		assert.Equal(t, exp.Pts, cha.Pts)
	}
}
func TestCountCribAction(t *testing.T) {
	g := model.Game{
		Crib: []model.Card{
			model.NewCardFromString(`2c`),
			model.NewCardFromString(`3c`),
			model.NewCardFromString(`4c`),
			model.NewCardFromString(`5c`),
		},
		CutCard: model.NewCardFromString(`10h`),
	}
	tests := []struct {
		desc string
		npc  model.PlayerID
		g    model.Game
		exp  model.PlayerAction
	}{{
		desc: `test dumb npc`,
		npc:  Dumb,
		exp: model.PlayerAction{
			ID:        Dumb,
			Overcomes: model.CountCrib,
			Action: model.CountCribAction{
				Pts: 8,
			}},
	}, {
		desc: `test simple npc`,
		npc:  Simple,
		exp: model.PlayerAction{
			ID:        Simple,
			Overcomes: model.CountCrib,
			Action: model.CountCribAction{
				Pts: 8,
			}},
	}, {
		desc: `test calculated npc`,
		npc:  Calc,
		exp: model.PlayerAction{
			ID:        Calc,
			Overcomes: model.CountCrib,
			Action: model.CountCribAction{
				Pts: 8,
			}},
	}}
	for _, tc := range tests {
		p := createPlayer(t, tc.npc)

		a := p.buildAction(tc.exp.Overcomes, g)
		assert.Equal(t, a.Overcomes, tc.exp.Overcomes)

		cca, ok := a.Action.(model.CountCribAction)
		assert.True(t, ok)
		exp, ok := tc.exp.Action.(model.CountCribAction)
		assert.True(t, ok)
		assert.Equal(t, exp.Pts, cca.Pts)
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
		g:     newGame(`dumbNPC`, 2, make([]model.Card, 0)),
		expGo: false,
	}, {
		desc:  `test simple npc`,
		npc:   `simpleNPC`,
		g:     newGame(`simpleNPC`, 2, make([]model.Card, 0)),
		expGo: false,
	}, {
		desc:  `test calculated npc`,
		npc:   `calculatedNPC`,
		g:     newGame(`calculatedNPC`, 2, make([]model.Card, 0)),
		expGo: false,
	}, {
		desc: `test go`,
		npc:  `dumbNPC`,
		g: newGame(`dumbNPC`, 2, []model.Card{
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
		g:         newGame(`dumbNPC`, 2, make([]model.Card, 0)),
		expNCards: 2,
	}, {
		desc:      `test simple npc`,
		npc:       `simpleNPC`,
		g:         newGame(`simpleNPC`, 2, make([]model.Card, 0)),
		expNCards: 2,
	}, {
		desc:      `test calculated npc`,
		npc:       `calculatedNPC`,
		g:         newGame(`calculatedNPC`, 2, make([]model.Card, 0)),
		expNCards: 2,
	}, {
		desc:      `test 3 player game`,
		npc:       `dumbNPC`,
		g:         newGame(`dumbNPC`, 3, make([]model.Card, 0)),
		expNCards: 1,
	}, {
		desc:      `test 4 player game`,
		npc:       `dumbNPC`,
		g:         newGame(`dumbNPC`, 4, make([]model.Card, 0)),
		expNCards: 1,
	}}
	for _, tc := range tests {
		p := createPlayer(t, tc.npc)

		for i := 0; i < 1; i++ {
			a := p.buildAction(model.CribCard, tc.g)
			assert.Equal(t, a.Overcomes, model.CribCard)

			bca, ok := a.Action.(model.BuildCribAction)
			assert.True(t, ok)
			assert.Len(t, bca.Cards, tc.expNCards, tc.desc)
		}
	}
}
