package network

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/joshprzybyszewski/cribbage/model"
)

func TestConvertToTeams(t *testing.T) {
	aliceID := model.PlayerID(`alice`)
	aliceName := `alice`
	bobID := model.PlayerID(`bob`)
	bobName := `bob`

	gID := model.NewGameID()

	testCases := []struct {
		desc  string
		input model.Game
		exp   []GetGameResponseTeam
	}{{
		desc: `undeclared teams`,
		input: model.Game{
			ID:            gID,
			CurrentScores: map[model.PlayerColor]int{},
			LagScores:     map[model.PlayerColor]int{},
			Players: []model.Player{{
				ID:   aliceID,
				Name: aliceName,
				Games: map[model.GameID]model.PlayerColor{
					gID: model.UnsetColor,
				},
			}, {
				ID:    bobID,
				Name:  bobName,
				Games: map[model.GameID]model.PlayerColor{},
			}},
		},
		exp: []GetGameResponseTeam{{
			Players: []Player{{
				ID:   aliceID,
				Name: aliceName,
			}},
			Color: ``,
		}, {
			Players: []Player{{
				ID:   bobID,
				Name: bobName,
			}},
			Color: ``,
		}},
	}, {
		desc: `one declared team`,
		input: model.Game{
			ID: gID,
			CurrentScores: map[model.PlayerColor]int{
				model.Red: 1,
			},
			LagScores: map[model.PlayerColor]int{
				model.Red: 0,
			},
			Players: []model.Player{{
				ID:   aliceID,
				Name: aliceName,
				Games: map[model.GameID]model.PlayerColor{
					gID: model.Red,
				},
			}, {
				ID:   bobID,
				Name: bobName,
				Games: map[model.GameID]model.PlayerColor{
					gID: model.UnsetColor,
				},
			}},
		},
		exp: []GetGameResponseTeam{{
			Players: []Player{{
				ID:   aliceID,
				Name: aliceName,
			}},
			Color:        `red`,
			CurrentScore: 1,
			LagScore:     0,
		}, {
			Players: []Player{{
				ID:   bobID,
				Name: bobName,
			}},
			Color:        ``,
			CurrentScore: 0,
			LagScore:     0,
		}},
	}, {
		desc: `two player game declared both team`,
		input: model.Game{
			ID: gID,
			CurrentScores: map[model.PlayerColor]int{
				model.Red:  1,
				model.Blue: 3,
			},
			LagScores: map[model.PlayerColor]int{
				model.Red:  0,
				model.Blue: 2,
			},
			Players: []model.Player{{
				ID:   aliceID,
				Name: aliceName,
				Games: map[model.GameID]model.PlayerColor{
					gID: model.Blue,
				},
			}, {
				ID:   bobID,
				Name: bobName,
				Games: map[model.GameID]model.PlayerColor{
					gID: model.Red,
				},
			}},
		},
		exp: []GetGameResponseTeam{{
			Players: []Player{{
				ID:   aliceID,
				Name: aliceName,
			}},
			Color:        `blue`,
			CurrentScore: 3,
			LagScore:     2,
		}, {
			Players: []Player{{
				ID:   bobID,
				Name: bobName,
			}},
			Color:        `red`,
			CurrentScore: 1,
			LagScore:     0,
		}},
	}}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			act := convertToTeams(tc.input)
			assert.Equal(t, tc.exp, act)
		})
	}
}
