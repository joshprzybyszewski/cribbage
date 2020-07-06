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
	erinID := model.PlayerID(`erin`)
	francisID := model.PlayerID(`francis`)

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
				Players: []model.Player{{
					ID:   chelseaID,
					Name: `chelsea`,
				}, {
					ID:   daveID,
					Name: `dave`,
				}},
				PlayerColors: map[model.PlayerID]model.PlayerColor{
					chelseaID: model.Red,
					daveID:    model.Blue,
				},
			},
			789: {
				Players: []model.Player{{
					ID:   erinID,
					Name: `erin`,
				}, {
					ID:   francisID,
					Name: `francis`,
				}},
				PlayerColors: map[model.PlayerID]model.PlayerColor{
					erinID:    model.Red,
					francisID: model.Blue,
				},
			},
		},
		expResp: GetActiveGamesForPlayerResponse{
			Player: Player{
				ID:   aliceID,
				Name: `alice`,
			},
			ActiveGames: map[model.GameID]ActiveGame{
				123: ActiveGame{
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
				},
				456: ActiveGame{
					PlayerNamesByID: map[model.PlayerID]string{
						chelseaID: `chelsea`,
						daveID:    `dave`,
					},
					PlayerColorsByID: map[model.PlayerID]string{
						chelseaID: `red`,
						daveID:    `blue`,
					},
					Created:  time.Time{},
					LastMove: time.Time{},
				},
				789: ActiveGame{
					PlayerNamesByID: map[model.PlayerID]string{
						erinID:    `erin`,
						francisID: `francis`,
					},
					PlayerColorsByID: map[model.PlayerID]string{
						erinID:    `red`,
						francisID: `blue`,
					},
					Created:  time.Time{},
					LastMove: time.Time{},
				},
			},
		},
	}}
	for _, tc := range tests {
		resp := ConvertToGetActiveGamesForPlayerResponse(tc.player, tc.inputGames)
		assert.Equal(t, tc.expResp, resp, tc.desc)
	}
}
