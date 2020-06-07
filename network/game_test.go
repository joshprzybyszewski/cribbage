package network

import (
	"testing"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/stretchr/testify/assert"
)

func handFromStrings(cs ...string) []model.Card {
	hand := make([]model.Card, len(cs))
	for i, c := range cs {
		hand[i] = model.NewCardFromString(c)
	}
	return hand
}

func TestNewGetGameResponse(t *testing.T) {
	tests := []struct {
		desc    string
		game    model.Game
		expResp GetGameResponse
	}{{
		desc: `shouldn't return hands or crib`,
		game: model.Game{
			ID: model.GameID(123456),
			Players: []model.Player{{
				ID:   `a`,
				Name: `a`,
			}, {
				ID:   `b`,
				Name: `b`,
			}},
			PlayerColors: map[model.PlayerID]model.PlayerColor{
				`a`: model.Blue,
				`b`: model.Red,
			},
			CurrentScores: map[model.PlayerColor]int{
				model.Blue: 11,
				model.Red:  22,
			},
			LagScores: map[model.PlayerColor]int{
				model.Blue: 10,
				model.Red:  20,
			},
			Phase: model.CribCounting,
			BlockingPlayers: map[model.PlayerID]model.Blocker{
				`b`: model.CountCrib,
			},
			CurrentDealer: `b`,
			Hands: map[model.PlayerID][]model.Card{
				`a`: handFromStrings(`ah`, `2h`, `3h`, `4h`),
				`b`: handFromStrings(`as`, `2s`, `3s`, `4s`),
			},
			Crib:    handFromStrings(`5h`, `6h`, `5s`, `6s`),
			CutCard: model.NewCardFromString(`5c`),
			PeggedCards: []model.PeggedCard{{
				Card:     model.NewCardFromString(`ah`),
				Action:   0,
				PlayerID: `a`,
			}, {
				Card:     model.NewCardFromString(`as`),
				Action:   1,
				PlayerID: `b`,
			}, {
				Card:     model.NewCardFromString(`2h`),
				Action:   2,
				PlayerID: `a`,
			}, {
				Card:     model.NewCardFromString(`2s`),
				Action:   3,
				PlayerID: `b`,
			}, {
				Card:     model.NewCardFromString(`3h`),
				Action:   4,
				PlayerID: `a`,
			}, {
				Card:     model.NewCardFromString(`3s`),
				Action:   5,
				PlayerID: `b`,
			}, {
				Card:     model.NewCardFromString(`4h`),
				Action:   6,
				PlayerID: `a`,
			}, {
				Card:     model.NewCardFromString(`4s`),
				Action:   7,
				PlayerID: `b`,
			}},
		},
		expResp: GetGameResponse{
			ID: model.GameID(123456),
			Players: []Player{{
				ID:   `a`,
				Name: `a`,
			}, {
				ID:   `b`,
				Name: `b`,
			}},
			PlayerColors: map[model.PlayerID]string{
				`a`: `blue`,
				`b`: `red`,
			},
			CurrentScores: map[string]int{
				`blue`: 11,
				`red`:  22,
			},
			LagScores: map[string]int{
				`blue`: 10,
				`red`:  20,
			},
			Phase: `CribCounting`,
			BlockingPlayers: map[model.PlayerID]string{
				`b`: `CountCrib`,
			},
			CurrentDealer: `b`,
			Hands:         nil,
			Crib:          nil,
			CutCard: Card{
				Suit:  `Clubs`,
				Value: 5,
				Name:  `5C`,
			},
			PeggedCards: []Card{{
				Suit:  `Hearts`,
				Value: 1,
				Name:  `AH`,
			}, {
				Suit:  `Spades`,
				Value: 1,
				Name:  `AS`,
			}, {
				Suit:  `Hearts`,
				Value: 2,
				Name:  `2H`,
			}, {
				Suit:  `Spades`,
				Value: 2,
				Name:  `2S`,
			}, {
				Suit:  `Hearts`,
				Value: 3,
				Name:  `3H`,
			}, {
				Suit:  `Spades`,
				Value: 3,
				Name:  `3S`,
			}, {
				Suit:  `Hearts`,
				Value: 4,
				Name:  `4H`,
			}, {
				Suit:  `Spades`,
				Value: 4,
				Name:  `4S`,
			}},
		},
	}}
	for _, tc := range tests {
		resp := NewGetGameResponse(tc.game)
		assert.Equal(t, tc.expResp, resp, tc.desc)
	}
}
