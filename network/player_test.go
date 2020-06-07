package network

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/joshprzybyszewski/cribbage/model"
)

func TestNewGetPlayerResponseFromModel(t *testing.T) {
	tests := []struct {
		desc    string
		player  model.Player
		expResp GetPlayerResponse
	}{{
		desc: ``,
		player: model.Player{
			ID:   `a`,
			Name: `aa`,
			Games: map[model.GameID]model.PlayerColor{
				123: model.Blue,
				456: model.Red,
				789: model.Green,
			},
		},
		expResp: GetPlayerResponse{
			Player: Player{
				ID:   `a`,
				Name: `aa`,
			},
			Games: map[model.GameID]string{
				123: `blue`,
				456: `red`,
				789: `green`,
			},
		},
	}}
	for _, tc := range tests {
		resp := NewGetPlayerResponseFromModel(tc.player)
		assert.Equal(t, tc.expResp, resp, tc.desc)
	}
}
func TestNewCreatePlayerResponseFromModel(t *testing.T) {
	tests := []struct {
		desc    string
		player  model.Player
		expResp CreatePlayerResponse
	}{{
		desc: ``,
		player: model.Player{
			ID:   `a`,
			Name: `aa`,
			Games: map[model.GameID]model.PlayerColor{
				123: model.Blue,
				456: model.Red,
				789: model.Green,
			},
		},
		expResp: CreatePlayerResponse{
			Player: Player{
				ID:   `a`,
				Name: `aa`,
			},
		},
	}}
	for _, tc := range tests {
		resp := NewCreatePlayerResponseFromModel(tc.player)
		assert.Equal(t, tc.expResp, resp, tc.desc)
	}
}
