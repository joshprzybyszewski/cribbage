package interaction

import (
	"testing"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/stretchr/testify/assert"
)

func Test_npcPlayer_buildDealAction(t *testing.T) {
	tests := []struct {
		desc string
		npc  NPCType
		g    model.Game
		exp  model.PlayerAction
	}{{
		desc: `test dumb npc`,
		npc:  Dumb,
	}, {
		desc: `test simple npc`,
		npc:  Simple,
	}, {
		desc: `test calculated npc`,
		npc:  Calculated,
	}}
	for _, tc := range tests {
		npc := npcPlayer{
			Type: tc.npc,
		}
		a := npc.buildAction(model.DealCards, tc.g)
		assert.Equal(t, a.Overcomes, model.DealCards)

		da, ok := a.Action.(model.DealAction)
		assert.True(t, ok)
		assert.Less(t, da.NumShuffles, 10)
		assert.GreaterOrEqual(t, da.NumShuffles, 1)
	}
}
func Test_npcPlayer_buildBuildCribAction(t *testing.T) {
	tests := []struct {
		desc      string
		npc       NPCType
		g         model.Game
		expNCards int
	}{{
		desc:      `test dumb npc`,
		npc:       Dumb,
		g:         model.Game{},
		expNCards: 2,
	}, {
		desc:      `test simple npc`,
		npc:       Simple,
		g:         model.Game{},
		expNCards: 2,
	}, {
		desc:      `test calculated npc`,
		npc:       Calculated,
		g:         model.Game{},
		expNCards: 2,
	}, {
		desc: `test 3 player game`,
		npc:  Calculated,
		g: model.Game{
			Players: make([]model.Player, 3),
		},
		expNCards: 1,
	}, {
		desc: `test 4 player game`,
		npc:  Calculated,
		g: model.Game{
			Players: make([]model.Player, 4),
		},
		expNCards: 1,
	}}
	for _, tc := range tests {
		npc := npcPlayer{
			Type: tc.npc,
		}
		a := npc.buildAction(model.CribCard, tc.g)
		assert.Equal(t, a.Overcomes, model.CribCard)

		bca, ok := a.Action.(model.BuildCribAction)
		assert.True(t, ok)
		assert.Len(t, bca.Cards, tc.expNCards)
	}
}

func TestNewNPCPlayer(t *testing.T) {
	tests := []struct {
		desc  string
		n     NPCType
		expID model.PlayerID
	}{
		{
			desc:  `test dumb npc`,
			n:     Dumb,
			expID: `dumbNPC`,
		},
		{
			desc:  `test simple npc`,
			n:     Simple,
			expID: `simpleNPC`,
		},
		{
			desc:  `test calculated npc`,
			n:     Calculated,
			expID: `calculatedNPC`,
		},
	}
	f := func(a model.PlayerAction) error {
		return nil
	}
	for _, tc := range tests {
		n := NewNPCPlayer(tc.n, f)
		assert.Equal(t, n.ID(), tc.expID)
	}
}
