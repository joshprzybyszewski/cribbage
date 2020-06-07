package network

import (
	"testing"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/stretchr/testify/assert"
)

func modelCardsFromStrings(cs ...string) []model.Card {
	hand := make([]model.Card, len(cs))
	for i, c := range cs {
		hand[i] = model.NewCardFromString(c)
	}
	return hand
}

func cardsFromStrings(cs ...string) []Card {
	hand := make([]Card, len(cs))
	for i, c := range cs {
		hand[i] = newCardFromModel(model.NewCardFromString(c))
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
				`a`: modelCardsFromStrings(`ah`, `2h`, `3h`, `4h`),
				`b`: modelCardsFromStrings(`as`, `2s`, `3s`, `4s`),
			},
			Crib:    modelCardsFromStrings(`5h`, `6h`, `5s`, `6s`),
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
			PeggedCards: cardsFromStrings(`AH`, `AS`, `2H`, `2S`, `3H`, `3S`, `4H`, `4S`),
		},
	}}
	for _, tc := range tests {
		resp := NewGetGameResponse(tc.game)
		assert.Equal(t, tc.expResp, resp, tc.desc)
	}
}
func TestNewGetGameResponseForPlayer(t *testing.T) {
	tests := []struct {
		desc    string
		player  model.PlayerID
		expErr  bool
		game    model.Game
		expResp GetGameResponse
	}{{
		desc:   `shouldn't return the other player's hand or the crib in early phase`,
		player: `a`,
		expErr: false,
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
			Phase: model.Pegging,
			BlockingPlayers: map[model.PlayerID]model.Blocker{
				`a`: model.PegCard,
			},
			CurrentDealer: `b`,
			Hands: map[model.PlayerID][]model.Card{
				`a`: modelCardsFromStrings(`ah`, `2h`, `3h`, `4h`),
				`b`: modelCardsFromStrings(`as`, `2s`, `3s`, `4s`),
			},
			Crib:    modelCardsFromStrings(`5h`, `6h`, `5s`, `6s`),
			CutCard: model.NewCardFromString(`5c`),
			PeggedCards: []model.PeggedCard{{
				Card:     model.NewCardFromString(`ah`),
				Action:   0,
				PlayerID: `a`,
			}, {
				Card:     model.NewCardFromString(`as`),
				Action:   1,
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
			Phase: `Pegging`,
			BlockingPlayers: map[model.PlayerID]string{
				`a`: `PegCard`,
			},
			CurrentDealer: `b`,
			Hands: map[model.PlayerID][]Card{
				`a`: cardsFromStrings(`AH`, `2H`, `3H`, `4H`),
			},
			Crib: nil,
			CutCard: Card{
				Suit:  `Clubs`,
				Value: 5,
				Name:  `5C`,
			},
			PeggedCards: cardsFromStrings(`AH`, `AS`),
		},
	}, {
		desc:   `shouldn't return the other player's hand or the crib in early phase`,
		player: `b`,
		expErr: false,
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
			Phase: model.Pegging,
			BlockingPlayers: map[model.PlayerID]model.Blocker{
				`a`: model.PegCard,
			},
			CurrentDealer: `b`,
			Hands: map[model.PlayerID][]model.Card{
				`a`: modelCardsFromStrings(`ah`, `2h`, `3h`, `4h`),
				`b`: modelCardsFromStrings(`as`, `2s`, `3s`, `4s`),
			},
			Crib:    modelCardsFromStrings(`5h`, `6h`, `5s`, `6s`),
			CutCard: model.NewCardFromString(`5c`),
			PeggedCards: []model.PeggedCard{{
				Card:     model.NewCardFromString(`ah`),
				Action:   0,
				PlayerID: `a`,
			}, {
				Card:     model.NewCardFromString(`as`),
				Action:   1,
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
			Phase: `Pegging`,
			BlockingPlayers: map[model.PlayerID]string{
				`a`: `PegCard`,
			},
			CurrentDealer: `b`,
			Hands: map[model.PlayerID][]Card{
				`b`: cardsFromStrings(`AS`, `2S`, `3S`, `4S`),
			},
			Crib: nil,
			CutCard: Card{
				Suit:  `Clubs`,
				Value: 5,
				Name:  `5C`,
			},
			PeggedCards: cardsFromStrings(`AH`, `AS`),
		},
	}, {
		desc:   `should return both hands but no crib after counting`,
		player: `b`,
		expErr: false,
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
			Phase: model.Counting,
			BlockingPlayers: map[model.PlayerID]model.Blocker{
				`a`: model.CountHand,
			},
			CurrentDealer: `b`,
			Hands: map[model.PlayerID][]model.Card{
				`a`: modelCardsFromStrings(`ah`, `2h`, `3h`, `4h`),
				`b`: modelCardsFromStrings(`as`, `2s`, `3s`, `4s`),
			},
			Crib:    modelCardsFromStrings(`5h`, `6h`, `5s`, `6s`),
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
			Phase: `Counting`,
			BlockingPlayers: map[model.PlayerID]string{
				`a`: `CountHand`,
			},
			CurrentDealer: `b`,
			Hands: map[model.PlayerID][]Card{
				`a`: cardsFromStrings(`AH`, `2H`, `3H`, `4H`),
				`b`: cardsFromStrings(`AS`, `2S`, `3S`, `4S`),
			},
			Crib: nil,
			CutCard: Card{
				Suit:  `Clubs`,
				Value: 5,
				Name:  `5C`,
			},
			PeggedCards: cardsFromStrings(`AH`, `AS`, `2H`, `2S`, `3H`, `3S`, `4H`, `4S`),
		},
	}, {
		desc:   `should return both hands and crib after counting crib`,
		player: `b`,
		expErr: false,
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
				`a`: modelCardsFromStrings(`ah`, `2h`, `3h`, `4h`),
				`b`: modelCardsFromStrings(`as`, `2s`, `3s`, `4s`),
			},
			Crib:    modelCardsFromStrings(`5h`, `6h`, `5s`, `6s`),
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
			Hands: map[model.PlayerID][]Card{
				`a`: cardsFromStrings(`AH`, `2H`, `3H`, `4H`),
				`b`: cardsFromStrings(`AS`, `2S`, `3S`, `4S`),
			},
			Crib: cardsFromStrings(`5H`, `6H`, `5S`, `6S`),
			CutCard: Card{
				Suit:  `Clubs`,
				Value: 5,
				Name:  `5C`,
			},
			PeggedCards: cardsFromStrings(`AH`, `AS`, `2H`, `2S`, `3H`, `3S`, `4H`, `4S`),
		},
	}, {
		desc:   `player not in game`,
		player: `c`,
		expErr: true,
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
			Phase: model.Pegging,
			BlockingPlayers: map[model.PlayerID]model.Blocker{
				`a`: model.PegCard,
			},
			CurrentDealer: `b`,
			Hands: map[model.PlayerID][]model.Card{
				`a`: modelCardsFromStrings(`ah`, `2h`, `3h`, `4h`),
				`b`: modelCardsFromStrings(`as`, `2s`, `3s`, `4s`),
			},
			Crib:    modelCardsFromStrings(`5h`, `6h`, `5s`, `6s`),
			CutCard: model.NewCardFromString(`5c`),
		},
		expResp: GetGameResponse{},
	}}
	for _, tc := range tests {
		resp, err := NewGetGameResponseForPlayer(tc.game, tc.player)

		if tc.expErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, tc.expResp, resp, tc.desc)
	}
}
