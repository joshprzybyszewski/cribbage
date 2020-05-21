package jsonutils

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshprzybyszewski/cribbage/model"
)

func addActionFor(g *model.Game, pID model.PlayerID, b model.Blocker, action interface{}) {
	pa := model.PlayerAction{
		GameID:    g.ID,
		ID:        pID,
		Overcomes: b,
		Action:    action,
	}
	g.Actions = append(g.Actions, pa)
}

func createGameWithActions() model.Game {
	p1 := model.PlayerID(`p1`)
	p2 := model.PlayerID(`p2`)
	g := model.Game{
		ID: model.GameID(1234),
		Players: []model.Player{{
			ID:   p1,
			Name: `p1`,
		}, {
			ID:   p2,
			Name: `p2`,
		}},
	}
	addActionFor(&g, p1, model.DealCards, model.DealAction{
		NumShuffles: 69,
	})
	addActionFor(&g, p1, model.CribCard, model.BuildCribAction{
		Cards: []model.Card{
			model.NewCardFromNumber(0),
			model.NewCardFromNumber(1),
		},
	})
	addActionFor(&g, p2, model.CribCard, model.BuildCribAction{
		Cards: []model.Card{
			model.NewCardFromNumber(2),
			model.NewCardFromNumber(3),
		},
	})
	return g
}

func TestBind(t *testing.T) {
	tests := []struct {
		name    string
		r       *http.Request
		obj     interface{}
		expGame model.Game
		wantErr bool
	}{{
		name:    `nil req`,
		r:       nil,
		obj:     nil,
		expGame: model.Game{},
		wantErr: true,
	}, {
		name: `nil req body`,
		r: &http.Request{
			Body: nil,
		},
		obj:     nil,
		expGame: model.Game{},
		wantErr: true,
	}}
	for _, tc := range tests {
		err := GameBinding.Bind(tc.r, tc.obj)
		if tc.wantErr {
			assert.Error(t, err)
			return
		}
		assert.NoError(t, err)
		assert.IsType(t, (*model.Game)(nil), tc.obj)
		g := tc.obj.(*model.Game)
		assert.Equal(t, tc.expGame, *g)
	}
}
func TestBindGame(t *testing.T) {
	tests := []struct {
		name    string
		obj     interface{}
		expGame model.Game
		expErr  bool
	}{{
		name:    `game with some actions`,
		obj:     &model.Game{},
		expGame: createGameWithActions(),
		expErr:  false,
	}, {
		name:    `try with not a model.Game`,
		obj:     model.Player{},
		expGame: model.Game{},
		expErr:  true,
	}}
	for _, tc := range tests {
		body, err := json.Marshal(tc.expGame)
		require.NoError(t, err, tc.name)
		err = bindGame(body, tc.obj)
		if tc.expErr {
			assert.Error(t, err, tc.name)
			return
		}
		assert.NoError(t, err, tc.name)
		game := tc.obj.(*model.Game)
		assert.Equal(t, tc.expGame, *game, tc.name)
	}
}
