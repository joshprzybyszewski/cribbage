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

	t2 := time.Now()
	t1 := t2.Add(-time.Minute)

	tests := []struct {
		desc       string
		player     model.Player
		inputGames map[model.GameID]model.Game
		expResp    GetActiveGamesForPlayerResponse
	}{{
		desc: `knows how to sort based on timestamp`,
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
				Actions: []model.PlayerAction{{
					TimeStamp: t1,
				}},
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
				Actions: []model.PlayerAction{{
					TimeStamp: t1,
				}, {
					TimeStamp: t2,
				}},
			},
		},
		expResp: GetActiveGamesForPlayerResponse{
			Player: Player{
				ID:   aliceID,
				Name: `alice`,
			},
			ActiveGames: []ActiveGame{{
				GameID: 789,
				Players: []ActiveGamePlayer{{
					ID:    aliceID,
					Name:  `alice`,
					Color: `red`,
				}, {
					ID:    daveID,
					Name:  `dave`,
					Color: `blue`,
				}},
				Created:  t1,
				LastMove: t2,
			}, {
				GameID: 456,
				Players: []ActiveGamePlayer{{
					ID:    aliceID,
					Name:  `alice`,
					Color: `red`,
				}, {
					ID:    chelseaID,
					Name:  `chelsea`,
					Color: `blue`,
				}},
				Created:  t1,
				LastMove: t1,
			}, {
				GameID: 123,
				Players: []ActiveGamePlayer{{
					ID:    aliceID,
					Name:  `alice`,
					Color: `red`,
				}, {
					ID:    bobID,
					Name:  `bob`,
					Color: `blue`,
				}},
				Created:  time.Time{},
				LastMove: time.Time{},
			}},
		},
	}, {
		desc: `knows how to filter games where alice isn't playing`,
		player: model.Player{
			ID:   aliceID,
			Name: `alice`,
			Games: map[model.GameID]model.PlayerColor{
				123: model.Blue,
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
					ID:   bobID,
					Name: `bob`,
				}, {
					ID:   chelseaID,
					Name: `chelsea`,
				}},
				PlayerColors: map[model.PlayerID]model.PlayerColor{
					bobID:     model.Red,
					chelseaID: model.Blue,
				},
			},
			789: {
				ID: 789,
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
		},
		expResp: GetActiveGamesForPlayerResponse{
			Player: Player{
				ID:   aliceID,
				Name: `alice`,
			},
			ActiveGames: []ActiveGame{{
				GameID: 123,
				Players: []ActiveGamePlayer{{
					ID:    aliceID,
					Name:  `alice`,
					Color: `red`,
				}, {
					ID:    bobID,
					Name:  `bob`,
					Color: `blue`,
				}},
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
