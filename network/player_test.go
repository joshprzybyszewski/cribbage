package network

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/joshprzybyszewski/cribbage/model"
)

func TestConvertToGetPlayerResponse(t *testing.T) {
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
		},
	}}
	for _, tc := range tests {
		resp := ConvertToGetPlayerResponse(tc.player)
		assert.Equal(t, tc.expResp, resp, tc.desc)
	}
}
func TestConvertToCreatePlayerResponse(t *testing.T) {
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
		resp := ConvertToCreatePlayerResponse(tc.player)
		assert.Equal(t, tc.expResp, resp, tc.desc)
	}
}

func TestConvertToGetActiveGamesForPlayerResponse(t *testing.T) {
	tests := []struct {
		desc       string
		player     model.Player
		inputGames map[model.GameID]model.Game
		expResp    GetActiveGamesForPlayerResponse
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
		inputGames: map[model.GameID]model.Game{
			123: {
				Players: []model.Player{{
					Name: `alice`,
				}, {
					Name: `bob`,
				}},
			},
			456: {
				Players: []model.Player{{
					Name: `chelsea`,
				}, {
					Name: `dave`,
				}},
			},
			789: {
				Players: []model.Player{{
					Name: `erin`,
				}, {
					Name: `francis`,
				}},
			},
		},
		expResp: GetActiveGamesForPlayerResponse{
			Player: Player{
				ID:   `a`,
				Name: `aa`,
			},
			ActiveGames: map[model.GameID]string{
				123: `alice, bob`,
				456: `chelsea, dave`,
				789: `erin, francis`,
			},
		},
	}}
	for _, tc := range tests {
		resp := ConvertToGetActiveGamesForPlayerResponse(tc.player, tc.inputGames)
		assert.Equal(t, tc.expResp, resp, tc.desc)
	}
}
