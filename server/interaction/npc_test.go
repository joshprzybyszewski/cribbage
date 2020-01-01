package interaction

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshprzybyszewski/cribbage/model"
)

func createPlayer(t *testing.T, pID model.PlayerID) *NPCPlayer {
	npc, err := NewNPCPlayer(pID, func(ctx context.Context, a model.PlayerAction) error {
		return nil
	})
	require.Nil(t, err)
	p, ok := npc.(*NPCPlayer)
	assert.True(t, ok)
	return p
}

func newGame(npcID model.PlayerID, nPlayers int, pegCards []model.Card) model.Game {
	players := make([]model.Player, nPlayers)
	for i := 0; i < nPlayers-1; i++ {
		id := model.PlayerID(fmt.Sprintf(`p%d`, i))
		players[i] = model.Player{ID: id}
	}
	players[len(players)-1] = model.Player{ID: npcID}

	hands := make(map[model.PlayerID][]model.Card)
	nCards := 6
	switch nPlayers {
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
			PlayerID: players[i%nPlayers].ID,
		})
	}
	return model.Game{
		ID:          5,
		Players:     players,
		Hands:       hands,
		PeggedCards: pegs,
	}
}

func TestDealAction(t *testing.T) {
	tests := []struct {
		desc string
		npc  model.PlayerID
	}{{
		desc: `test dumb npc`,
		npc:  Dumb,
	}, {
		desc: `test simple npc`,
		npc:  Simple,
	}, {
		desc: `test calculated npc`,
		npc:  Calc,
	}}
	for _, tc := range tests {
		p := createPlayer(t, tc.npc)

		a, err := p.buildAction(model.DealCards, model.Game{})
		assert.Nil(t, err)
		assert.Equal(t, a.Overcomes, model.DealCards)

		da, ok := a.Action.(model.DealAction)
		assert.True(t, ok)
		assert.LessOrEqual(t, da.NumShuffles, 10)
		assert.GreaterOrEqual(t, da.NumShuffles, 1)
	}
}
func TestCutAction(t *testing.T) {
	tests := []struct {
		desc string
		npc  model.PlayerID
	}{{
		desc: `test dumb npc`,
		npc:  Dumb,
	}, {
		desc: `test simple npc`,
		npc:  Simple,
	}, {
		desc: `test calculated npc`,
		npc:  Calc,
	}}
	for _, tc := range tests {
		p := createPlayer(t, tc.npc)

		a, err := p.buildAction(model.CutCard, model.Game{})
		assert.Nil(t, err)
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

		a, err := p.buildAction(model.CountHand, g)
		assert.Nil(t, err)
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

		a, err := p.buildAction(tc.exp.Overcomes, g)
		assert.Nil(t, err)
		assert.Equal(t, a.Overcomes, tc.exp.Overcomes)

		cca, ok := a.Action.(model.CountCribAction)
		assert.True(t, ok)
		exp, ok := tc.exp.Action.(model.CountCribAction)
		assert.True(t, ok)
		assert.Equal(t, exp.Pts, cca.Pts)
	}
}

func TestPegTwice(t *testing.T) {
	tests := []struct {
		desc  string
		npc   model.PlayerID
		g     model.Game
		expGo bool
	}{{
		desc:  `test dumb npc`,
		npc:   Dumb,
		g:     newGame(Dumb, 2, make([]model.Card, 0)),
		expGo: false,
	}, {
		desc:  `test simple npc`,
		npc:   Simple,
		g:     newGame(Simple, 2, make([]model.Card, 0)),
		expGo: false,
	}, {
		desc:  `test calculated npc`,
		npc:   Calc,
		g:     newGame(Calc, 2, make([]model.Card, 0)),
		expGo: false,
	}}
	for _, tc := range tests {
		p := createPlayer(t, tc.npc)
		a, err := p.buildAction(model.PegCard, tc.g)
		assert.NoError(t, err)
		assert.Equal(t, a.Overcomes, model.PegCard)

		pa, ok := a.Action.(model.PegAction)
		assert.True(t, ok)
		if !pa.SayGo {
			tc.g.PeggedCards = append(tc.g.PeggedCards, model.PeggedCard{
				Card:     pa.Card,
				PlayerID: a.ID,
			})
		}

		a, err = p.buildAction(model.PegCard, tc.g)
		assert.NoError(t, err)
		assert.Equal(t, a.Overcomes, model.PegCard)

		pa, ok = a.Action.(model.PegAction)
		assert.True(t, ok)
		c := model.PeggedCard{
			Card:     pa.Card,
			PlayerID: a.ID,
		}
		assert.NotContains(t, tc.g.PeggedCards, c)
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
		npc:   Dumb,
		g:     newGame(Dumb, 2, make([]model.Card, 0)),
		expGo: false,
	}, {
		desc:  `test simple npc`,
		npc:   Simple,
		g:     newGame(Simple, 2, make([]model.Card, 0)),
		expGo: false,
	}, {
		desc:  `test calculated npc`,
		npc:   Calc,
		g:     newGame(Calc, 2, make([]model.Card, 0)),
		expGo: false,
	}, {
		desc: `test dumb go`,
		npc:  Dumb,
		g: newGame(Dumb, 2, []model.Card{
			model.NewCardFromString(`10c`),
			model.NewCardFromString(`10s`),
			model.NewCardFromString(`10h`),
		}),
		expGo: true,
	}, {
		desc: `test simple go`,
		npc:  Simple,
		g: newGame(Simple, 2, []model.Card{
			model.NewCardFromString(`10c`),
			model.NewCardFromString(`10s`),
			model.NewCardFromString(`10h`),
		}),
		expGo: true,
	}, {
		desc: `test calculated go`,
		npc:  Calc,
		g: newGame(Calc, 2, []model.Card{
			model.NewCardFromString(`10c`),
			model.NewCardFromString(`10s`),
			model.NewCardFromString(`10h`),
		}),
		expGo: true,
	}}
	for _, tc := range tests {
		p := createPlayer(t, tc.npc)

		for i := 0; i < 10; i++ {
			a, err := p.buildAction(model.PegCard, tc.g)
			assert.Nil(t, err)
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
}
func TestBuildCribAction(t *testing.T) {
	tests := []struct {
		desc      string
		npc       model.PlayerID
		isDealer  bool
		g         model.Game
		expNCards int
	}{{
		desc:      `test dumb npc`,
		npc:       Dumb,
		isDealer:  false,
		g:         newGame(Dumb, 2, make([]model.Card, 0)),
		expNCards: 2,
	}, {
		desc:      `test simple npc, not dealer`,
		npc:       Simple,
		isDealer:  false,
		g:         newGame(Simple, 2, make([]model.Card, 0)),
		expNCards: 2,
	}, {
		desc:      `test simple npc, dealer`,
		npc:       Simple,
		isDealer:  true,
		g:         newGame(Simple, 2, make([]model.Card, 0)),
		expNCards: 2,
	}, {
		desc:      `test calculated npc, not dealer`,
		npc:       Calc,
		isDealer:  false,
		g:         newGame(Calc, 2, make([]model.Card, 0)),
		expNCards: 2,
	}, {
		desc:      `test calculated npc, dealer`,
		npc:       Calc,
		isDealer:  true,
		g:         newGame(Calc, 2, make([]model.Card, 0)),
		expNCards: 2,
	}, {
		desc:      `test 3 player game`,
		npc:       Dumb,
		isDealer:  false,
		g:         newGame(Dumb, 3, make([]model.Card, 0)),
		expNCards: 1,
	}, {
		desc:      `test 4 player game`,
		npc:       Dumb,
		isDealer:  false,
		g:         newGame(Dumb, 4, make([]model.Card, 0)),
		expNCards: 1,
	}}
	for _, tc := range tests {
		p := createPlayer(t, tc.npc)
		if tc.isDealer {
			tc.g.CurrentDealer = tc.npc
		}

		for i := 0; i < 5; i++ {
			a, err := p.buildAction(model.CribCard, tc.g)
			assert.Nil(t, err)
			assert.Equal(t, a.Overcomes, model.CribCard)

			bca, ok := a.Action.(model.BuildCribAction)
			assert.True(t, ok)
			assert.Len(t, bca.Cards, tc.expNCards, tc.desc)
		}
	}
}

func TestNotifyBlocking(t *testing.T) {
	tests := []struct {
		desc string
		npc  model.PlayerID
	}{{
		desc: `test dumb NPC`,
		npc:  Dumb,
	}, {
		desc: `test simple NPC`,
		npc:  Simple,
	}, {
		desc: `test calculated NPC`,
		npc:  Calc,
	}}

	for _, tc := range tests {
		cb := func(ctx context.Context, a model.PlayerAction) error {
			da, ok := a.Action.(model.DealAction)
			assert.True(t, ok)
			assert.GreaterOrEqual(t, da.NumShuffles, 1)
			assert.LessOrEqual(t, da.NumShuffles, 10)
			return nil
		}
		p, err := NewNPCPlayer(tc.npc, cb)
		require.Nil(t, err)
		err = p.NotifyBlocking(model.DealCards, model.Game{}, ``)
		assert.Nil(t, err)
	}
}
func TestNotifyMessage(t *testing.T) {
	tests := []struct {
		desc string
		npc  model.PlayerID
	}{{
		desc: `test dumb NPC`,
		npc:  Dumb,
	}, {
		desc: `test simple NPC`,
		npc:  Simple,
	}, {
		desc: `test calculated NPC`,
		npc:  Calc,
	}}

	for _, tc := range tests {
		p, err := NewNPCPlayer(tc.npc, func(ctx context.Context, a model.PlayerAction) error {
			return nil
		})
		require.Nil(t, err)
		err = p.NotifyMessage(model.Game{}, ``)
		assert.Nil(t, err)
	}
}
func TestNotifyScoreUpdate(t *testing.T) {
	tests := []struct {
		desc string
		npc  model.PlayerID
	}{{
		desc: `test dumb NPC`,
		npc:  Dumb,
	}, {
		desc: `test simple NPC`,
		npc:  Simple,
	}, {
		desc: `test calculated NPC`,
		npc:  Calc,
	}}

	for _, tc := range tests {
		p, err := NewNPCPlayer(tc.npc, func(ctx context.Context, a model.PlayerAction) error {
			return nil
		})
		require.Nil(t, err)
		err = p.NotifyScoreUpdate(model.Game{}, ``)
		assert.Nil(t, err)
	}
}
func TestNewNPCPlayer(t *testing.T) {
	tests := []struct {
		desc      string
		npc       model.PlayerID
		expErr    bool
		expPlayer interface{}
	}{{
		desc:      `test dumb NPC`,
		npc:       Dumb,
		expErr:    false,
		expPlayer: &dumbNPC{},
	}, {
		desc:      `test simple NPC`,
		npc:       Simple,
		expErr:    false,
		expPlayer: &simpleNPC{},
	}, {
		desc:      `test calculated NPC`,
		npc:       Calc,
		expErr:    false,
		expPlayer: &calculatedNPC{},
	}, {
		desc:      `test unsupported type`,
		npc:       `unsupported`,
		expErr:    true,
		expPlayer: nil,
	}}

	for _, tc := range tests {
		p, err := NewNPCPlayer(tc.npc, func(ctx context.Context, a model.PlayerAction) error {
			return nil
		})
		if tc.expErr {
			assert.Error(t, err)
			assert.Nil(t, p)
		} else {
			n, ok := p.(*NPCPlayer)
			assert.True(t, ok)
			assert.NoError(t, err)
			assert.Equal(t, tc.expPlayer, n.player)
		}
	}
}
