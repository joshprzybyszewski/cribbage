package network

import (
	"testing"
	"time"

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
	aliceID := model.PlayerID(`alice`)
	bobID := model.PlayerID(`bob`)
	chelseaID := model.PlayerID(`chelsea`)
	daveID := model.PlayerID(`dave`)

	tests := []struct {
		desc       string
		player     model.Player
		inputGames map[model.GameID]model.Game
		expResp    GetActiveGamesForPlayerResponse
	}{{
		desc: ``,
		player: model.Player{
			ID:   aliceID,
			Name: `alice`,
			Games: map[model.GameID]model.PlayerColor{
				123: model.Blue,
				456: model.Red,
				789: model.Green,
			},
		},
		inputGames: map[model.GameID]model.Game{
			123: {
				ID: 123,
				Players: []model.Player{{
					ID:   aliceID,
					Name: `alice`,
				}, {
					ID:   bobID,
					Name: `bob`,
				}},
				PlayerColors: map[model.PlayerID]model.PlayerColor{
					aliceID: model.Red,
					bobID:   model.Blue,
				},
			},
			456: {
				ID: 456,
				Players: []model.Player{{
					ID:   aliceID,
					Name: `alice`,
				}, {
					ID:   chelseaID,
					Name: `chelsea`,
				}},
				PlayerColors: map[model.PlayerID]model.PlayerColor{
					aliceID:   model.Red,
					chelseaID: model.Blue,
				},
			},
			789: {
				ID: 789,
				Players: []model.Player{{
					ID:   aliceID,
					Name: `alice`,
				}, {
					ID:   daveID,
					Name: `dave`,
				}},
				PlayerColors: map[model.PlayerID]model.PlayerColor{
					aliceID: model.Red,
					daveID:  model.Blue,
				},
			},
		},
		expResp: GetActiveGamesForPlayerResponse{
			Player: Player{
				ID:   aliceID,
				Name: `alice`,
			},
			ActiveGames: []ActiveGame{{
				GameID: 123,
				PlayerNamesByID: map[model.PlayerID]string{
					aliceID: `alice`,
					bobID:   `bob`,
				},
				PlayerColorsByID: map[model.PlayerID]string{
					aliceID: `red`,
					bobID:   `blue`,
				},
				Created:  time.Time{},
				LastMove: time.Time{},
			}, {
				GameID: 456,
				PlayerNamesByID: map[model.PlayerID]string{
					aliceID:   `alice`,
					chelseaID: `chelsea`,
				},
				PlayerColorsByID: map[model.PlayerID]string{
					aliceID:   `red`,
					chelseaID: `blue`,
				},
				Created:  time.Time{},
				LastMove: time.Time{},
			}, {
				GameID: 789,
				PlayerNamesByID: map[model.PlayerID]string{
					aliceID: `alice`,
					daveID:  `dave`,
				},
				PlayerColorsByID: map[model.PlayerID]string{
					aliceID: `red`,
					daveID:  `blue`,
				},
				Created:  time.Time{},
				LastMove: time.Time{},
			}},
		},
	}}
	for _, tc := range tests {
		resp := ConvertToGetActiveGamesForPlayerResponse(tc.player, tc.inputGames)
		assert.Equal(t, tc.expResp, resp, tc.desc)
	}
}
