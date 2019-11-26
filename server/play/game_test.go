package play

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
)

func setup() (a, b model.Player, am, bm *interaction.Mock, pAPIs map[model.PlayerID]interaction.Player) {
	alice := model.Player{
		ID:   model.PlayerID(1),
		Name: `alice`,
	}
	bob := model.Player{
		ID:   model.PlayerID(2),
		Name: `bob`,
	}
	aAPI := &interaction.Mock{}
	bAPI := &interaction.Mock{}
	abAPIs := map[model.PlayerID]interaction.Player{
		alice.ID: aAPI,
		bob.ID:   bAPI,
	}
	return alice, bob, aAPI, bAPI, abAPIs
}

func TestHandleAction_Deal(t *testing.T) {
	alice, bob, aliceAPI, bobAPI, abAPIs := setup()

	g := model.Game{
		ID:              model.GameID(5),
		Players:         []model.Player{alice, bob},
		Deck:            model.NewDeck(),
		BlockingPlayers: map[model.PlayerID]model.Blocker{alice.ID: model.DealCards},
		CurrentDealer:   alice.ID,
		PlayerColors:    map[model.PlayerID]model.PlayerColor{alice.ID: model.Blue, bob.ID: model.Red},
		CurrentScores:   map[model.PlayerColor]int{model.Blue: 0, model.Red: 0},
		LagScores:       map[model.PlayerColor]int{model.Blue: 0, model.Red: 0},
		Phase:           model.Deal,
		Hands:           make(map[model.PlayerID][]model.Card, 2),
		CutCard:         model.Card{},
		Crib:            make([]model.Card, 4),
		PeggedCards:     make([]model.PeggedCard, 0, 8),
		CanResetPeg:     false,
	}
	action := model.PlayerAction{
		GameID:    g.ID,
		ID:        alice.ID,
		Overcomes: model.DealCards,
		Action: model.DealAction{
			NumShuffles: 50,
		},
	}
	aliceAPI.On(`NotifyMessage`, mock.AnythingOfType(`model.Game`), mock.MatchedBy(func(s string) bool { return strings.HasPrefix(s, `Received Hand `) })).Return(nil).Once()
	aliceAPI.On(`NotifyBlocking`, model.CribCard, mock.AnythingOfType(`model.Game`), `needs to cut 2 cards`).Return(nil).Once()
	bobAPI.On(`NotifyMessage`, mock.AnythingOfType(`model.Game`), mock.MatchedBy(func(s string) bool { return strings.HasPrefix(s, `Received Hand `) })).Return(nil).Once()
	bobAPI.On(`NotifyBlocking`, model.CribCard, mock.AnythingOfType(`model.Game`), `needs to cut 2 cards`).Return(nil).Once()

	err := HandleAction(&g, action, abAPIs)
	assert.Nil(t, err)
	assert.Equal(t, model.BuildCrib, g.Phase)
	assert.Equal(t, g.NumActions, 1)
	// now the game is blocked by both players needing to build the crib
	require.Len(t, g.BlockingPlayers, 2)
	assert.Contains(t, g.BlockingPlayers, alice.ID)
	assert.Contains(t, g.BlockingPlayers, bob.ID)
	// the players should have 6 card hands
	assert.Len(t, g.Hands[alice.ID], 6)
	assert.Len(t, g.Hands[bob.ID], 6)
	// assert that entering the build crib phase has cleared out the crib
	assert.Empty(t, g.Crib)

	aliceAPI.AssertExpectations(t)
	bobAPI.AssertExpectations(t)
}

func TestHandleAction_Crib(t *testing.T) {
	alice, bob, aliceAPI, bobAPI, abAPIs := setup()

	g := model.Game{
		ID:              model.GameID(5),
		NumActions:      1,
		Players:         []model.Player{alice, bob},
		Deck:            model.NewDeck(),
		BlockingPlayers: map[model.PlayerID]model.Blocker{alice.ID: model.CribCard, bob.ID: model.CribCard},
		CurrentDealer:   alice.ID,
		PlayerColors:    map[model.PlayerID]model.PlayerColor{alice.ID: model.Blue, bob.ID: model.Red},
		CurrentScores:   map[model.PlayerColor]int{model.Blue: 0, model.Red: 0},
		LagScores:       map[model.PlayerColor]int{model.Blue: 0, model.Red: 0},
		Phase:           model.BuildCrib,
		Hands: map[model.PlayerID][]model.Card{
			alice.ID: []model.Card{
				model.NewCardFromString(`1s`),
				model.NewCardFromString(`2s`),
				model.NewCardFromString(`3s`),
				model.NewCardFromString(`4s`),
				model.NewCardFromString(`5s`),
				model.NewCardFromString(`6s`),
			},
			bob.ID: []model.Card{
				model.NewCardFromString(`1c`),
				model.NewCardFromString(`2c`),
				model.NewCardFromString(`3c`),
				model.NewCardFromString(`4c`),
				model.NewCardFromString(`5c`),
				model.NewCardFromString(`6c`),
			},
		},
		CutCard:     model.NewCardFromString(`KH`),
		Crib:        make([]model.Card, 0, 4),
		PeggedCards: make([]model.PeggedCard, 0, 8),
		CanResetPeg: false,
	}

	action := model.PlayerAction{
		GameID:    g.ID,
		ID:        bob.ID,
		Overcomes: model.CribCard,
		Action: model.BuildCribAction{
			Cards: []model.Card{
				model.NewCardFromString(`1c`),
				model.NewCardFromString(`2c`),
			},
		},
	}

	err := HandleAction(&g, action, abAPIs)
	assert.Nil(t, err)
	assert.Equal(t, model.BuildCrib, g.Phase)
	assert.Equal(t, g.NumActions, 2)
	// now the game is blocked by alice needing to submit a crib card
	require.Len(t, g.BlockingPlayers, 1)
	assert.Contains(t, g.BlockingPlayers, alice.ID)
	assert.NotContains(t, g.BlockingPlayers, bob.ID)
	// alice should have 6 card hands, bob only has 4
	assert.Len(t, g.Hands[alice.ID], 6)
	assert.Len(t, g.Hands[bob.ID], 4)
	assert.Len(t, g.Crib, 2)

	action = model.PlayerAction{
		GameID:    g.ID,
		ID:        alice.ID,
		Overcomes: model.CribCard,
		Action: model.BuildCribAction{
			Cards: []model.Card{
				model.NewCardFromString(`1s`),
				model.NewCardFromString(`2s`),
			},
		},
	}
	bobAPI.On(`NotifyBlocking`, model.CutCard, mock.AnythingOfType(`model.Game`), ``).Return(nil).Once()

	err = HandleAction(&g, action, abAPIs)
	assert.Nil(t, err)
	assert.Equal(t, model.Cut, g.Phase)
	assert.Equal(t, g.NumActions, 3)
	// now the game has moved on to cutting for bob
	require.Len(t, g.BlockingPlayers, 1)
	assert.NotContains(t, g.BlockingPlayers, alice.ID)
	assert.Contains(t, g.BlockingPlayers, bob.ID)
	// the players hand should all be developed, and the crib too
	assert.Len(t, g.Hands[alice.ID], 4)
	assert.Contains(t, g.Hands[alice.ID], model.NewCardFromString(`3s`))
	assert.Contains(t, g.Hands[alice.ID], model.NewCardFromString(`4s`))
	assert.Contains(t, g.Hands[alice.ID], model.NewCardFromString(`5s`))
	assert.Contains(t, g.Hands[alice.ID], model.NewCardFromString(`6s`))
	assert.Len(t, g.Hands[bob.ID], 4)
	assert.Contains(t, g.Hands[bob.ID], model.NewCardFromString(`3c`))
	assert.Contains(t, g.Hands[bob.ID], model.NewCardFromString(`4c`))
	assert.Contains(t, g.Hands[bob.ID], model.NewCardFromString(`5c`))
	assert.Contains(t, g.Hands[bob.ID], model.NewCardFromString(`6c`))
	assert.Len(t, g.Crib, 4)
	assert.Contains(t, g.Crib, model.NewCardFromString(`1s`))
	assert.Contains(t, g.Crib, model.NewCardFromString(`2s`))
	assert.Contains(t, g.Crib, model.NewCardFromString(`1c`))
	assert.Contains(t, g.Crib, model.NewCardFromString(`2c`))
	// verify that entering the cutting phase clears out the cut card until it _is_ cut
	assert.Equal(t, model.Card{}, g.CutCard)

	aliceAPI.AssertExpectations(t)
	bobAPI.AssertExpectations(t)
}

func TestHandleAction_Pegging(t *testing.T) {
	alice, bob, aliceAPI, bobAPI, abAPIs := setup()

	g := model.Game{
		ID: model.GameID(5),
		// TODO
		NumActions:      1,
		Players:         []model.Player{alice, bob},
		Deck:            model.NewDeck(),
		BlockingPlayers: map[model.PlayerID]model.Blocker{bob.ID: model.PegCard},
		CurrentDealer:   alice.ID,
		PlayerColors:    map[model.PlayerID]model.PlayerColor{alice.ID: model.Blue, bob.ID: model.Red},
		CurrentScores:   map[model.PlayerColor]int{model.Blue: 0, model.Red: 0},
		LagScores:       map[model.PlayerColor]int{model.Blue: 0, model.Red: 0},
		Phase:           model.Pegging,
		// cards are chosen so we can test a 31 and a go
		Hands: map[model.PlayerID][]model.Card{
			alice.ID: []model.Card{
				model.NewCardFromString(`7s`),
				model.NewCardFromString(`8s`),
				model.NewCardFromString(`10s`),
				model.NewCardFromString(`js`),
			},
			bob.ID: []model.Card{
				model.NewCardFromString(`7c`),
				model.NewCardFromString(`9c`),
				model.NewCardFromString(`10c`),
				model.NewCardFromString(`jc`),
			},
		},
		CutCard: model.NewCardFromString(`KH`),
		Crib: []model.Card{
			model.NewCardFromString(`as`),
			model.NewCardFromString(`ah`),
			model.NewCardFromString(`ac`),
			model.NewCardFromString(`ad`),
		},
		PeggedCards: make([]model.PeggedCard, 0, 8),
		CanResetPeg: false,
	}

	action := model.PlayerAction{
		GameID:    g.ID,
		ID:        bob.ID,
		Overcomes: model.PegCard,
		Action: model.PegAction{
			Card:  g.Hands[bob.ID][0],
			SayGo: false,
		},
	}
	aliceAPI.On(`NotifyBlocking`, model.PegCard, mock.AnythingOfType(`model.Game`), ``).Return(nil).Once()
	err := HandleAction(&g, action, abAPIs)
	assert.Nil(t, err)
	assert.Len(t, g.PeggedCards, 1)
	assert.Contains(t, g.PeggedCards, model.PeggedCard{Card: model.NewCardFromString(`7c`), PlayerID: bob.ID})
	assert.Equal(t, g.CurrentPeg(), 7)

	action = model.PlayerAction{
		GameID:    g.ID,
		ID:        alice.ID,
		Overcomes: model.PegCard,
		Action: model.PegAction{
			Card:  g.Hands[alice.ID][0],
			SayGo: false,
		},
	}
	// alice and bob are going to get notified because alice scores a 31 and a run of 3
	aliceAPI.On(`NotifyScoreUpdate`, mock.AnythingOfType(`model.Game`), []string{`pegging`}).Return(nil).Once()
	bobAPI.On(`NotifyScoreUpdate`, mock.AnythingOfType(`model.Game`), []string{`pegging`}).Return(nil).Once()
	bobAPI.On(`NotifyBlocking`, model.PegCard, mock.AnythingOfType(`model.Game`), ``).Return(nil).Once()
	err = HandleAction(&g, action, abAPIs)
	assert.Nil(t, err)
	assert.Len(t, g.PeggedCards, 2)
	assert.Contains(t, g.PeggedCards, model.PeggedCard{Card: model.NewCardFromString(`7s`), PlayerID: alice.ID})
	assert.Equal(t, g.CurrentPeg(), 14)

	action = model.PlayerAction{
		GameID:    g.ID,
		ID:        bob.ID,
		Overcomes: model.PegCard,
		Action: model.PegAction{
			Card:  g.Hands[bob.ID][1],
			SayGo: false,
		},
	}

	aliceAPI.On(`NotifyBlocking`, model.PegCard, mock.AnythingOfType(`model.Game`), ``).Return(nil).Once()
	err = HandleAction(&g, action, abAPIs)
	assert.Nil(t, err)
	assert.Len(t, g.PeggedCards, 3)
	assert.Contains(t, g.PeggedCards, model.PeggedCard{Card: model.NewCardFromString(`9c`), PlayerID: bob.ID})
	assert.Equal(t, g.CurrentPeg(), 23)

	action = model.PlayerAction{
		GameID:    g.ID,
		ID:        alice.ID,
		Overcomes: model.PegCard,
		Action: model.PegAction{
			Card:  g.Hands[alice.ID][1],
			SayGo: false,
		},
	}
	// alice and bob are going to get notified because alice scores a 31 and a run of 3
	aliceAPI.On(`NotifyScoreUpdate`, mock.AnythingOfType(`model.Game`), []string{`pegging`}).Return(nil).Once()
	bobAPI.On(`NotifyScoreUpdate`, mock.AnythingOfType(`model.Game`), []string{`pegging`}).Return(nil).Once()
	bobAPI.On(`NotifyBlocking`, model.PegCard, mock.AnythingOfType(`model.Game`), ``).Return(nil).Once()
	err = HandleAction(&g, action, abAPIs)
	assert.Nil(t, err)
	assert.Len(t, g.PeggedCards, 4)
	assert.Contains(t, g.PeggedCards, model.PeggedCard{Card: model.NewCardFromString(`8s`), PlayerID: alice.ID})
	assert.Equal(t, g.CurrentPeg(), 31)
	assert.Equal(t, g.CurrentScores[g.PlayerColors[alice.ID]], 7)

	action = model.PlayerAction{
		GameID:    g.ID,
		ID:        bob.ID,
		Overcomes: model.PegCard,
		Action: model.PegAction{
			Card:  g.Hands[bob.ID][2],
			SayGo: false,
		},
	}

	aliceAPI.On(`NotifyBlocking`, model.PegCard, mock.AnythingOfType(`model.Game`), ``).Return(nil).Once()
	err = HandleAction(&g, action, abAPIs)
	assert.Nil(t, err)
	assert.Len(t, g.PeggedCards, 5)
	assert.Contains(t, g.PeggedCards, model.PeggedCard{Card: model.NewCardFromString(`10c`), PlayerID: bob.ID})
	assert.Equal(t, g.CurrentPeg(), 10)

	action = model.PlayerAction{
		GameID:    g.ID,
		ID:        alice.ID,
		Overcomes: model.PegCard,
		Action: model.PegAction{
			Card:  g.Hands[alice.ID][2],
			SayGo: false,
		},
	}
	// alice and bob are going to get notified because alice scores a pair
	aliceAPI.On(`NotifyScoreUpdate`, mock.AnythingOfType(`model.Game`), []string{`pegging`}).Return(nil).Once()
	bobAPI.On(`NotifyScoreUpdate`, mock.AnythingOfType(`model.Game`), []string{`pegging`}).Return(nil).Once()
	bobAPI.On(`NotifyBlocking`, model.PegCard, mock.AnythingOfType(`model.Game`), ``).Return(nil).Once()
	err = HandleAction(&g, action, abAPIs)
	assert.Nil(t, err)
	assert.Len(t, g.PeggedCards, 6)
	assert.Contains(t, g.PeggedCards, model.PeggedCard{Card: model.NewCardFromString(`10s`), PlayerID: alice.ID})
	assert.Equal(t, g.CurrentPeg(), 20)
	assert.Equal(t, g.CurrentScores[g.PlayerColors[alice.ID]], 9)

	action = model.PlayerAction{
		GameID:    g.ID,
		ID:        bob.ID,
		Overcomes: model.PegCard,
		Action: model.PegAction{
			Card:  g.Hands[bob.ID][3],
			SayGo: false,
		},
	}

	aliceAPI.On(`NotifyBlocking`, model.PegCard, mock.AnythingOfType(`model.Game`), ``).Return(nil).Once()
	err = HandleAction(&g, action, abAPIs)
	assert.Nil(t, err)
	assert.Len(t, g.PeggedCards, 7)
	assert.Contains(t, g.PeggedCards, model.PeggedCard{Card: model.NewCardFromString(`jc`), PlayerID: bob.ID})
	assert.Equal(t, g.CurrentPeg(), 30)

	action = model.PlayerAction{
		GameID:    g.ID,
		ID:        alice.ID,
		Overcomes: model.PegCard,
		Action: model.PegAction{
			Card:  g.Hands[alice.ID][3],
			SayGo: false,
		},
	}
	// alice tries to peg a jack but she needs to say go
	aliceAPI.On(`NotifyBlocking`, model.PegCard, mock.AnythingOfType(`model.Game`), `Cannot peg card with this value`).Return(nil).Once()
	err = HandleAction(&g, action, abAPIs)
	assert.Nil(t, err)
	assert.Len(t, g.PeggedCards, 7)

	// alice says go
	action = model.PlayerAction{
		GameID:    g.ID,
		ID:        alice.ID,
		Overcomes: model.PegCard,
		Action: model.PegAction{
			Card:  model.Card{},
			SayGo: true,
		},
	}
	bobAPI.On(`NotifyBlocking`, model.PegCard, mock.AnythingOfType(`model.Game`), ``).Return(nil).Once()
	err = HandleAction(&g, action, abAPIs)
	assert.Nil(t, err)
	assert.Len(t, g.PeggedCards, 7)

	// bob says go
	action = model.PlayerAction{
		GameID:    g.ID,
		ID:        bob.ID,
		Overcomes: model.PegCard,
		Action: model.PegAction{
			Card:  model.Card{},
			SayGo: true,
		},
	}
	// bob scores a go
	aliceAPI.On(`NotifyScoreUpdate`, mock.AnythingOfType(`model.Game`), []string{`the go`}).Return(nil).Once()
	bobAPI.On(`NotifyScoreUpdate`, mock.AnythingOfType(`model.Game`), []string{`the go`}).Return(nil).Once()
	aliceAPI.On(`NotifyBlocking`, model.PegCard, mock.AnythingOfType(`model.Game`), ``).Return(nil).Once()
	err = HandleAction(&g, action, abAPIs)
	assert.Nil(t, err)
	assert.Len(t, g.PeggedCards, 7)
	assert.Equal(t, g.CurrentScores[g.PlayerColors[bob.ID]], 1)

	// make sure we're still at 30 after all of that
	require.Equal(t, g.CurrentPeg(), 30)
	action = model.PlayerAction{
		GameID:    g.ID,
		ID:        alice.ID,
		Overcomes: model.PegCard,
		Action: model.PegAction{
			Card:  g.Hands[alice.ID][3],
			SayGo: false,
		},
	}
	// alice scores last card
	aliceAPI.On(`NotifyScoreUpdate`, mock.AnythingOfType(`model.Game`), []string{`last card`}).Return(nil).Once()
	bobAPI.On(`NotifyScoreUpdate`, mock.AnythingOfType(`model.Game`), []string{`last card`}).Return(nil).Once()
	// bob will be up to count his hand
	bobAPI.On(`NotifyBlocking`, model.CountHand, mock.AnythingOfType(`model.Game`), ``).Return(nil).Once()
	err = HandleAction(&g, action, abAPIs)
	assert.Nil(t, err)
	assert.Len(t, g.PeggedCards, 8)
	assert.Contains(t, g.PeggedCards, model.PeggedCard{Card: model.NewCardFromString(`js`), PlayerID: alice.ID})
	assert.Equal(t, g.CurrentPeg(), 10)
	assert.Equal(t, g.CurrentScores[g.PlayerColors[alice.ID]], 10)

	// we have moved on to counting hands, and bob is up
	assert.Equal(t, model.Counting, g.Phase)
	assert.Contains(t, g.BlockingPlayers, bob.ID)
	assert.NotContains(t, g.BlockingPlayers, alice.ID)

	aliceAPI.AssertExpectations(t)
	bobAPI.AssertExpectations(t)
}

func TestHandleAction_Counting(t *testing.T) {
	alice, bob, aliceAPI, bobAPI, abAPIs := setup()

	g := model.Game{
		ID: model.GameID(5),
		// TODO
		NumActions:      1,
		Players:         []model.Player{alice, bob},
		Deck:            model.NewDeck(),
		BlockingPlayers: map[model.PlayerID]model.Blocker{bob.ID: model.CountHand},
		CurrentDealer:   alice.ID,
		PlayerColors:    map[model.PlayerID]model.PlayerColor{alice.ID: model.Blue, bob.ID: model.Red},
		CurrentScores:   map[model.PlayerColor]int{model.Blue: 0, model.Red: 0},
		LagScores:       map[model.PlayerColor]int{model.Blue: 0, model.Red: 0},
		Phase:           model.Counting,
		Hands: map[model.PlayerID][]model.Card{
			alice.ID: []model.Card{
				model.NewCardFromString(`7s`),
				model.NewCardFromString(`8s`),
				model.NewCardFromString(`9s`),
				model.NewCardFromString(`10s`),
			},
			bob.ID: []model.Card{
				model.NewCardFromString(`7c`),
				model.NewCardFromString(`8c`),
				model.NewCardFromString(`9c`),
				model.NewCardFromString(`10c`),
			},
		},
		CutCard:     model.NewCardFromString(`7h`),
		Crib:        make([]model.Card, 4),
		PeggedCards: make([]model.PeggedCard, 0, 8),
		CanResetPeg: false,
	}

	action := model.PlayerAction{
		GameID:    g.ID,
		ID:        bob.ID,
		Overcomes: model.CountHand,
		Action: model.CountHandAction{
			Pts: 18,
		},
	}
	bobAPI.On(`NotifyScoreUpdate`, mock.AnythingOfType(`model.Game`), []string{`hand (7♥︎: 7♣︎, 8♣︎, 9♣︎, 10♣︎)`}).Return(nil).Once()
	aliceAPI.On(`NotifyScoreUpdate`, mock.AnythingOfType(`model.Game`), []string{`hand (7♥︎: 7♣︎, 8♣︎, 9♣︎, 10♣︎)`}).Return(nil).Once()
	aliceAPI.On(`NotifyBlocking`, model.CountHand, mock.AnythingOfType(`model.Game`), ``).Return(nil).Once()
	err := HandleAction(&g, action, abAPIs)
	assert.Nil(t, err)
	assert.Equal(t, 18, g.CurrentScores[g.PlayerColors[bob.ID]])
	assert.Contains(t, g.BlockingPlayers, alice.ID)
	assert.NotContains(t, g.BlockingPlayers, bob.ID)

	action = model.PlayerAction{
		GameID:    g.ID,
		ID:        alice.ID,
		Overcomes: model.CountHand,
		Action: model.CountHandAction{
			Pts: 18,
		},
	}
	bobAPI.On(`NotifyScoreUpdate`, mock.AnythingOfType(`model.Game`), []string{`hand (7♥︎: 7♠︎, 8♠︎, 9♠︎, 10♠︎)`}).Return(nil).Once()
	aliceAPI.On(`NotifyScoreUpdate`, mock.AnythingOfType(`model.Game`), []string{`hand (7♥︎: 7♠︎, 8♠︎, 9♠︎, 10♠︎)`}).Return(nil).Once()
	aliceAPI.On(`NotifyBlocking`, model.CountCrib, mock.AnythingOfType(`model.Game`), ``).Return(nil).Once()
	err = HandleAction(&g, action, abAPIs)
	assert.Nil(t, err)
	assert.Equal(t, 18, g.CurrentScores[g.PlayerColors[alice.ID]])

	// counting is done - we've moved onto counting the crib and alice needs to do that
	assert.Equal(t, model.CribCounting, g.Phase)
	assert.Contains(t, g.BlockingPlayers, alice.ID)
	assert.NotContains(t, g.BlockingPlayers, bob.ID)

	aliceAPI.AssertExpectations(t)
	bobAPI.AssertExpectations(t)
}
