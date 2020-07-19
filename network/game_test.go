package network

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/joshprzybyszewski/cribbage/model"
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
		if c == `` {
			hand[i] = invalidCard
		} else {
			hand[i] = convertToCard(model.NewCardFromString(c))
		}
	}
	return hand
}

func newPeggedCard(c string, pID model.PlayerID) PeggedCard {
	return PeggedCard{
		Card:   convertToCard(model.NewCardFromString(c)),
		Player: pID,
	}
}

func TestConvertToCreateGameResponse(t *testing.T) {
	tests := []struct {
		desc    string
		game    model.Game
		expResp CreateGameResponse
	}{{
		desc: ``,
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
				model.Blue: 0,
				model.Red:  0,
			},
			LagScores: map[model.PlayerColor]int{
				model.Red:  0,
				model.Blue: 0,
			},
			Phase: model.Deal,
			BlockingPlayers: map[model.PlayerID]model.Blocker{
				`b`: model.DealCards,
			},
			CurrentDealer: `b`,
			Hands:         map[model.PlayerID][]model.Card{},
			Crib:          []model.Card{},
			CutCard:       model.Card{},
			PeggedCards:   []model.PeggedCard{},
		},
		expResp: CreateGameResponse{
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
			BlockingPlayers: map[model.PlayerID]string{
				`b`: `DealCards`,
			},
			CurrentDealer: `b`,
		},
	}}
	for _, tc := range tests {
		resp := ConvertToCreateGameResponse(tc.game)
		assert.Equal(t, tc.expResp, resp, tc.desc)
	}
}

func TestConvertToGetGameResponse(t *testing.T) {
	aliceID := model.PlayerID(`alice`)
	bobID := model.PlayerID(`bob`)

	gID := model.GameID(123456)

	tests := []struct {
		desc    string
		game    model.Game
		expResp GetGameResponse
	}{{
		desc: `shouldn't return hands or crib`,
		game: model.Game{
			ID: gID,
			Players: []model.Player{{
				ID:   aliceID,
				Name: `alicia`,
				Games: map[model.GameID]model.PlayerColor{
					gID: model.Blue,
				},
			}, {
				ID:   bobID,
				Name: `bobbette`,
				Games: map[model.GameID]model.PlayerColor{
					gID: model.Red,
				},
			}},
			PlayerColors: map[model.PlayerID]model.PlayerColor{
				aliceID: model.Blue,
				bobID:   model.Red,
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
				bobID: model.CountCrib,
			},
			CurrentDealer: bobID,
			Hands: map[model.PlayerID][]model.Card{
				aliceID: modelCardsFromStrings(`ah`, `2h`, `3h`, `4h`),
				bobID:   modelCardsFromStrings(`as`, `2s`, `3s`, `4s`),
			},
			Crib:    modelCardsFromStrings(`5h`, `6h`, `5s`, `6s`),
			CutCard: model.NewCardFromString(`5c`),
			PeggedCards: []model.PeggedCard{{
				Card:     model.NewCardFromString(`ah`),
				Action:   0,
				PlayerID: aliceID,
			}, {
				Card:     model.NewCardFromString(`as`),
				Action:   1,
				PlayerID: bobID,
			}, {
				Card:     model.NewCardFromString(`2h`),
				Action:   2,
				PlayerID: aliceID,
			}, {
				Card:     model.NewCardFromString(`2s`),
				Action:   3,
				PlayerID: bobID,
			}, {
				Card:     model.NewCardFromString(`3h`),
				Action:   4,
				PlayerID: aliceID,
			}, {
				Card:     model.NewCardFromString(`3s`),
				Action:   5,
				PlayerID: bobID,
			}, {
				Card:     model.NewCardFromString(`4h`),
				Action:   6,
				PlayerID: aliceID,
			}, {
				Card:     model.NewCardFromString(`4s`),
				Action:   7,
				PlayerID: bobID,
			}},
		},
		expResp: GetGameResponse{
			ID: model.GameID(123456),
			Teams: []GetGameResponseTeam{{
				Color:        `blue`,
				CurrentScore: 11,
				LagScore:     10,
				Players: []Player{{
					ID:   aliceID,
					Name: `alicia`,
				}},
			}, {
				Color:        `red`,
				CurrentScore: 22,
				LagScore:     20,
				Players: []Player{{
					ID:   bobID,
					Name: `bobbette`,
				}},
			}},
			Phase: `CribCounting`,
			BlockingPlayers: map[model.PlayerID]string{
				bobID: `CountCrib`,
			},
			CurrentDealer: bobID,
			Hands:         nil,
			Crib:          cardsFromStrings(`5h`, `6h`, `5s`, `6s`),
			CutCard: Card{
				Suit:  `Clubs`,
				Value: 5,
				Name:  `5C`,
			},
			PeggedCards: []PeggedCard{
				newPeggedCard(`AH`, aliceID),
				newPeggedCard(`AS`, bobID),
				newPeggedCard(`2H`, aliceID),
				newPeggedCard(`2S`, bobID),
				newPeggedCard(`3H`, aliceID),
				newPeggedCard(`3S`, bobID),
				newPeggedCard(`4H`, aliceID),
				newPeggedCard(`4S`, bobID),
			},
		},
	}}
	for _, tc := range tests {
		resp := ConvertToGetGameResponse(tc.game)
		assert.Equal(t, tc.expResp, resp, tc.desc)
	}
}

func TestConvertToGetGameResponseForPlayer(t *testing.T) {
	aliceID := model.PlayerID(`alice`)
	bobID := model.PlayerID(`bob`)

	gID := model.GameID(123456)

	tests := []struct {
		desc    string
		player  model.PlayerID
		expErr  bool
		game    model.Game
		expResp GetGameResponse
	}{{
		desc:   `should only return the cards which have been revealed to alice`,
		player: aliceID,
		expErr: false,
		game: model.Game{
			ID: gID,
			Players: []model.Player{{
				ID:   aliceID,
				Name: `alice`,
				Games: map[model.GameID]model.PlayerColor{
					gID: model.Blue,
				},
			}, {
				ID:   bobID,
				Name: `bob`,
				Games: map[model.GameID]model.PlayerColor{
					gID: model.Red,
				},
			}},
			PlayerColors: map[model.PlayerID]model.PlayerColor{
				aliceID: model.Blue,
				bobID:   model.Red,
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
				aliceID: model.PegCard,
			},
			CurrentDealer: bobID,
			Hands: map[model.PlayerID][]model.Card{
				aliceID: modelCardsFromStrings(`ah`, `2h`, `3h`, `4h`),
				bobID:   modelCardsFromStrings(`as`, `2s`, `3s`, `4s`),
			},
			Crib:    modelCardsFromStrings(`5h`, `6h`, `5s`, `6s`),
			CutCard: model.NewCardFromString(`5c`),
			PeggedCards: []model.PeggedCard{{
				Card:     model.NewCardFromString(`ah`),
				Action:   0,
				PlayerID: aliceID,
			}, {
				Card:     model.NewCardFromString(`as`),
				Action:   1,
				PlayerID: bobID,
			}},
		},
		expResp: GetGameResponse{
			ID: model.GameID(123456),
			Teams: []GetGameResponseTeam{{
				Color:        `blue`,
				CurrentScore: 11,
				LagScore:     10,
				Players: []Player{{
					ID:   aliceID,
					Name: `alice`,
				}},
			}, {
				Color:        `red`,
				CurrentScore: 22,
				LagScore:     20,
				Players: []Player{{
					ID:   bobID,
					Name: `bob`,
				}},
			}},
			Phase: `Pegging`,
			BlockingPlayers: map[model.PlayerID]string{
				aliceID: `PegCard`,
			},
			CurrentDealer: bobID,
			Hands: map[model.PlayerID][]Card{
				aliceID: cardsFromStrings(`AH`, `2H`, `3H`, `4H`),
				bobID:   cardsFromStrings(`AS`, ``, ``, ``),
			},
			Crib: cardsFromStrings(``, ``, ``, ``),
			CutCard: Card{
				Suit:  `Clubs`,
				Value: 5,
				Name:  `5C`,
			},
			PeggedCards: []PeggedCard{
				newPeggedCard(`AH`, aliceID),
				newPeggedCard(`AS`, bobID),
			},
			CurrentPeg: 2,
		},
	}, {
		desc:   `should only return the cards which have been revealed to player b`,
		player: `b`,
		expErr: false,
		game: model.Game{
			ID: gID,
			Players: []model.Player{{
				ID:   aliceID,
				Name: `alice`,
				Games: map[model.GameID]model.PlayerColor{
					gID: model.Blue,
				},
			}, {
				ID:   bobID,
				Name: `bob`,
				Games: map[model.GameID]model.PlayerColor{
					gID: model.Red,
				},
			}},
			PlayerColors: map[model.PlayerID]model.PlayerColor{
				aliceID: model.Blue,
				bobID:   model.Red,
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
				aliceID: model.PegCard,
			},
			CurrentDealer: bobID,
			Hands: map[model.PlayerID][]model.Card{
				aliceID: modelCardsFromStrings(`ah`, `2h`, `3h`, `4h`),
				bobID:   modelCardsFromStrings(`as`, `2s`, `3s`, `4s`),
			},
			Crib:    modelCardsFromStrings(`5h`, `6h`, `5s`, `6s`),
			CutCard: model.NewCardFromString(`5c`),
			PeggedCards: []model.PeggedCard{{
				Card:     model.NewCardFromString(`ah`),
				Action:   0,
				PlayerID: aliceID,
			}, {
				Card:     model.NewCardFromString(`as`),
				Action:   1,
				PlayerID: bobID,
			}},
		},
		expResp: GetGameResponse{
			ID: model.GameID(123456),
			Teams: []GetGameResponseTeam{{
				Color:        `blue`,
				CurrentScore: 11,
				LagScore:     10,
				Players: []Player{{
					ID:   aliceID,
					Name: `alice`,
				}},
			}, {
				Color:        `red`,
				CurrentScore: 22,
				LagScore:     20,
				Players: []Player{{
					ID:   bobID,
					Name: `bob`,
				}},
			}},
			Phase: `Pegging`,
			BlockingPlayers: map[model.PlayerID]string{
				aliceID: `PegCard`,
			},
			CurrentDealer: bobID,
			Hands: map[model.PlayerID][]Card{
				aliceID: cardsFromStrings(`AH`, ``, ``, ``),
				bobID:   cardsFromStrings(`AS`, `2S`, `3S`, `4S`),
			},
			CurrentPeg: 2,
			Crib:       cardsFromStrings(``, ``, ``, ``),
			CutCard: Card{
				Suit:  `Clubs`,
				Value: 5,
				Name:  `5C`,
			},
			PeggedCards: []PeggedCard{
				newPeggedCard(`AH`, aliceID),
				newPeggedCard(`AS`, bobID),
			},
		},
	}, {
		desc:   `should return both hands but no crib after counting`,
		player: `b`,
		expErr: false,
		game: model.Game{
			ID: gID,
			Players: []model.Player{{
				ID:   aliceID,
				Name: `anne`,
				Games: map[model.GameID]model.PlayerColor{
					gID: model.Blue,
				},
			}, {
				ID:   bobID,
				Name: `bryan`,
				Games: map[model.GameID]model.PlayerColor{
					gID: model.Red,
				},
			}},
			PlayerColors: map[model.PlayerID]model.PlayerColor{
				aliceID: model.Blue,
				bobID:   model.Red,
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
				aliceID: model.CountHand,
			},
			CurrentDealer: bobID,
			Hands: map[model.PlayerID][]model.Card{
				aliceID: modelCardsFromStrings(`ah`, `2h`, `3h`, `4h`),
				bobID:   modelCardsFromStrings(`as`, `2s`, `3s`, `4s`),
			},
			Crib:    modelCardsFromStrings(`5h`, `6h`, `5s`, `6s`),
			CutCard: model.NewCardFromString(`5c`),
			PeggedCards: []model.PeggedCard{{
				Card:     model.NewCardFromString(`ah`),
				Action:   0,
				PlayerID: aliceID,
			}, {
				Card:     model.NewCardFromString(`as`),
				Action:   1,
				PlayerID: bobID,
			}, {
				Card:     model.NewCardFromString(`2h`),
				Action:   2,
				PlayerID: aliceID,
			}, {
				Card:     model.NewCardFromString(`2s`),
				Action:   3,
				PlayerID: bobID,
			}, {
				Card:     model.NewCardFromString(`3h`),
				Action:   4,
				PlayerID: aliceID,
			}, {
				Card:     model.NewCardFromString(`3s`),
				Action:   5,
				PlayerID: bobID,
			}, {
				Card:     model.NewCardFromString(`4h`),
				Action:   6,
				PlayerID: aliceID,
			}, {
				Card:     model.NewCardFromString(`4s`),
				Action:   7,
				PlayerID: bobID,
			}},
		},
		expResp: GetGameResponse{
			ID: model.GameID(123456),
			Teams: []GetGameResponseTeam{{
				Color:        `blue`,
				CurrentScore: 11,
				LagScore:     10,
				Players: []Player{{
					ID:   aliceID,
					Name: `anne`,
				}},
			}, {
				Color:        `red`,
				CurrentScore: 22,
				LagScore:     20,
				Players: []Player{{
					ID:   bobID,
					Name: `bryan`,
				}},
			}},
			Phase: `Counting`,
			BlockingPlayers: map[model.PlayerID]string{
				aliceID: `CountHand`,
			},
			CurrentDealer: bobID,
			Hands: map[model.PlayerID][]Card{
				aliceID: cardsFromStrings(`AH`, `2H`, `3H`, `4H`),
				bobID:   cardsFromStrings(`AS`, `2S`, `3S`, `4S`),
			},
			CurrentPeg: 0,
			Crib:       cardsFromStrings(``, ``, ``, ``),
			CutCard: Card{
				Suit:  `Clubs`,
				Value: 5,
				Name:  `5C`,
			},
			PeggedCards: []PeggedCard{
				newPeggedCard(`AH`, aliceID),
				newPeggedCard(`AS`, bobID),
				newPeggedCard(`2H`, aliceID),
				newPeggedCard(`2S`, bobID),
				newPeggedCard(`3H`, aliceID),
				newPeggedCard(`3S`, bobID),
				newPeggedCard(`4H`, aliceID),
				newPeggedCard(`4S`, bobID),
			},
		},
	}, {
		desc:   `should return both hands and crib after counting crib`,
		player: `b`,
		expErr: false,
		game: model.Game{
			ID: gID,
			Players: []model.Player{{
				ID:   aliceID,
				Name: `alishia`,
				Games: map[model.GameID]model.PlayerColor{
					gID: model.Blue,
				},
			}, {
				ID:   bobID,
				Name: `robert`,
				Games: map[model.GameID]model.PlayerColor{
					gID: model.Red,
				},
			}},
			PlayerColors: map[model.PlayerID]model.PlayerColor{
				aliceID: model.Blue,
				bobID:   model.Red,
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
				bobID: model.CountCrib,
			},
			CurrentDealer: bobID,
			Hands: map[model.PlayerID][]model.Card{
				aliceID: modelCardsFromStrings(`ah`, `2h`, `3h`, `4h`),
				bobID:   modelCardsFromStrings(`as`, `2s`, `3s`, `4s`),
			},
			Crib:    modelCardsFromStrings(`5h`, `6h`, `5s`, `6s`),
			CutCard: model.NewCardFromString(`5c`),
			PeggedCards: []model.PeggedCard{{
				Card:     model.NewCardFromString(`ah`),
				Action:   0,
				PlayerID: aliceID,
			}, {
				Card:     model.NewCardFromString(`as`),
				Action:   1,
				PlayerID: bobID,
			}, {
				Card:     model.NewCardFromString(`2h`),
				Action:   2,
				PlayerID: aliceID,
			}, {
				Card:     model.NewCardFromString(`2s`),
				Action:   3,
				PlayerID: bobID,
			}, {
				Card:     model.NewCardFromString(`3h`),
				Action:   4,
				PlayerID: aliceID,
			}, {
				Card:     model.NewCardFromString(`3s`),
				Action:   5,
				PlayerID: bobID,
			}, {
				Card:     model.NewCardFromString(`4h`),
				Action:   6,
				PlayerID: aliceID,
			}, {
				Card:     model.NewCardFromString(`4s`),
				Action:   7,
				PlayerID: bobID,
			}},
		},
		expResp: GetGameResponse{
			ID: model.GameID(123456),
			Teams: []GetGameResponseTeam{{
				Color:        `blue`,
				CurrentScore: 11,
				LagScore:     10,
				Players: []Player{{
					ID:   aliceID,
					Name: `alishia`,
				}},
			}, {
				Color:        `red`,
				CurrentScore: 22,
				LagScore:     20,
				Players: []Player{{
					ID:   bobID,
					Name: `robert`,
				}},
			}},
			Phase: `CribCounting`,
			BlockingPlayers: map[model.PlayerID]string{
				bobID: `CountCrib`,
			},
			CurrentDealer: bobID,
			Hands: map[model.PlayerID][]Card{
				aliceID: cardsFromStrings(`AH`, `2H`, `3H`, `4H`),
				bobID:   cardsFromStrings(`AS`, `2S`, `3S`, `4S`),
			},
			Crib: cardsFromStrings(`5H`, `6H`, `5S`, `6S`),
			CutCard: Card{
				Suit:  `Clubs`,
				Value: 5,
				Name:  `5C`,
			},
			PeggedCards: []PeggedCard{
				newPeggedCard(`AH`, aliceID),
				newPeggedCard(`AS`, bobID),
				newPeggedCard(`2H`, aliceID),
				newPeggedCard(`2S`, bobID),
				newPeggedCard(`3H`, aliceID),
				newPeggedCard(`3S`, bobID),
				newPeggedCard(`4H`, aliceID),
				newPeggedCard(`4S`, bobID),
			},
		},
	}, {
		desc:   `player not in game`,
		player: `c`,
		expErr: true,
		game: model.Game{
			ID: gID,
			Players: []model.Player{{
				ID:   aliceID,
				Name: `alicia`,
				Games: map[model.GameID]model.PlayerColor{
					gID: model.Blue,
				},
			}, {
				ID:   bobID,
				Name: `bobb`,
				Games: map[model.GameID]model.PlayerColor{
					gID: model.Red,
				},
			}},
			PlayerColors: map[model.PlayerID]model.PlayerColor{
				aliceID: model.Blue,
				bobID:   model.Red,
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
				aliceID: model.PegCard,
			},
			CurrentDealer: bobID,
			Hands: map[model.PlayerID][]model.Card{
				aliceID: modelCardsFromStrings(`ah`, `2h`, `3h`, `4h`),
				bobID:   modelCardsFromStrings(`as`, `2s`, `3s`, `4s`),
			},
			Crib:    modelCardsFromStrings(`5h`, `6h`, `5s`, `6s`),
			CutCard: model.NewCardFromString(`5c`),
		},
		expResp: GetGameResponse{},
	}}
	for _, tc := range tests {
		resp, err := ConvertToGetGameResponseForPlayer(tc.game, tc.player)

		if tc.expErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, tc.expResp, resp, tc.desc)
		if tc.expErr {
			continue
		}

		mg2 := ConvertFromGetGameResponse(resp)
		if mg2.Phase >= model.CribCounting {
			// by the time the crib comes around, the network and the model
			// games should be identical
			assert.Equal(t, mg2, tc.game, tc.desc)
			continue
		}
		assert.NotEqual(t, mg2, tc.game, tc.desc)

		assert.Len(t, mg2.Crib, len(tc.game.Crib), tc.desc)
		for pID, hand := range mg2.Hands {
			if pID == tc.player {
				assert.Equal(t, tc.game.Hands[pID], hand)
			} else {
				for _, c := range hand {
					wasPegged := c == model.InvalidCard
					for _, pc := range tc.game.PeggedCards {
						if pc.Card == c {
							wasPegged = true
							break
						}
					}
					assert.True(t, wasPegged, `expected revealed card (%v) to have been pegged, but wasn't. (%s)`, c, tc.desc)
				}
			}
		}

		// Make the hands field nil on both of the games. _Now_ they should be equal
		mg2.Hands = nil
		tc.game.Hands = nil
		for i := range tc.game.Crib {
			tc.game.Crib[i] = model.InvalidCard
		}
		assert.Equal(t, tc.game, mg2, tc.desc)
	}
}
