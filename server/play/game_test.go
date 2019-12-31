package play

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/utils/testutils"
)

func TestHandleAction_InvalidINputs(t *testing.T) {
	alice, bob, _, _, abAPIs := testutils.AliceAndBob()

	g := model.Game{
		ID:              model.GameID(5),
		Players:         []model.Player{alice, bob},
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
	}
	action := model.PlayerAction{
		GameID:    model.GameID(8),
		ID:        alice.ID,
		Overcomes: model.DealCards,
		Action: model.DealAction{
			NumShuffles: 50,
		},
	}

	err := HandleAction(&g, action, abAPIs)
	assert.Equal(t, ErrActionNotForGame, err)

	action.GameID = g.ID
	action.ID = model.PlayerID(`dne`)

	err = HandleAction(&g, action, abAPIs)
	assert.Equal(t, ErrPlayerNotInGame, err)

	action.ID = alice.ID
	g.CurrentScores[model.Blue] = 121

	err = HandleAction(&g, action, abAPIs)
	assert.Equal(t, ErrGameAlreadyOver, err)
}

func TestHandleAction_Deal(t *testing.T) {
	alice, bob, aliceAPI, bobAPI, abAPIs := testutils.AliceAndBob()

	g := model.Game{
		ID:              model.GameID(5),
		Players:         []model.Player{alice, bob},
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
	assert.Equal(t, g.NumActions(), 1)
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
	alice, bob, aliceAPI, bobAPI, abAPIs := testutils.AliceAndBob()

	g := model.Game{
		ID:              model.GameID(5),
		Players:         []model.Player{alice, bob},
		BlockingPlayers: map[model.PlayerID]model.Blocker{alice.ID: model.CribCard, bob.ID: model.CribCard},
		CurrentDealer:   alice.ID,
		PlayerColors:    map[model.PlayerID]model.PlayerColor{alice.ID: model.Blue, bob.ID: model.Red},
		CurrentScores:   map[model.PlayerColor]int{model.Blue: 0, model.Red: 0},
		LagScores:       map[model.PlayerColor]int{model.Blue: 0, model.Red: 0},
		Phase:           model.BuildCrib,
		Hands: map[model.PlayerID][]model.Card{
			alice.ID: {
				model.NewCardFromString(`1s`),
				model.NewCardFromString(`2s`),
				model.NewCardFromString(`3s`),
				model.NewCardFromString(`4s`),
				model.NewCardFromString(`5s`),
				model.NewCardFromString(`6s`),
			},
			bob.ID: {
				model.NewCardFromString(`1c`),
				model.NewCardFromString(`2c`),
				model.NewCardFromString(`3c`),
				model.NewCardFromString(`4c`),
				model.NewCardFromString(`5c`),
				model.NewCardFromString(`6c`),
			},
		},
		CutCard:     model.Card{}, //NewCardFromString(`KH`),
		Crib:        make([]model.Card, 0, 4),
		PeggedCards: make([]model.PeggedCard, 0, 8),
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
	assert.Equal(t, g.NumActions(), 1)
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
	assert.Equal(t, g.NumActions(), 2)
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

func TestHandleAction_Cut(t *testing.T) {
	alice, bob, aliceAPI, bobAPI, abAPIs := testutils.AliceAndBob()

	g := model.Game{
		ID:              model.GameID(5),
		Players:         []model.Player{alice, bob},
		BlockingPlayers: map[model.PlayerID]model.Blocker{bob.ID: model.CutCard},
		CurrentDealer:   alice.ID,
		PlayerColors:    map[model.PlayerID]model.PlayerColor{alice.ID: model.Blue, bob.ID: model.Red},
		CurrentScores:   map[model.PlayerColor]int{model.Blue: 0, model.Red: 0},
		LagScores:       map[model.PlayerColor]int{model.Blue: 0, model.Red: 0},
		Phase:           model.Cut,
		Hands: map[model.PlayerID][]model.Card{
			alice.ID: {
				model.NewCardFromString(`1s`),
				model.NewCardFromString(`2s`),
				model.NewCardFromString(`3s`),
				model.NewCardFromString(`4s`),
			},
			bob.ID: {
				model.NewCardFromString(`1c`),
				model.NewCardFromString(`2c`),
				model.NewCardFromString(`3c`),
				model.NewCardFromString(`4c`),
			},
		},
		CutCard: model.Card{},
		Crib: []model.Card{
			model.NewCardFromString(`6s`),
			model.NewCardFromString(`6c`),
			model.NewCardFromString(`6d`),
			model.NewCardFromString(`6h`),
		},
		PeggedCards: make([]model.PeggedCard, 0, 8),
	}

	action := model.PlayerAction{
		GameID:    g.ID,
		ID:        bob.ID,
		Overcomes: model.CutCard,
		Action: model.CutDeckAction{
			Percentage: 0.314,
		},
	}

	aliceAPI.On(`NotifyMessage`, mock.AnythingOfType(`model.Game`), mock.Anything).Return(nil).Once()
	bobAPI.On(`NotifyMessage`, mock.AnythingOfType(`model.Game`), mock.Anything).Return(nil).Once()
	bobAPI.On(`NotifyBlocking`, model.PegCard, mock.AnythingOfType(`model.Game`), `please peg a card`).Return(nil).Once()

	err := HandleAction(&g, action, abAPIs)
	assert.Nil(t, err)
	assert.Equal(t, model.Pegging, g.Phase)
	assert.Equal(t, g.NumActions(), 1)
	// now the game has moved on to pegging for bob
	require.Len(t, g.BlockingPlayers, 1)
	assert.NotContains(t, g.BlockingPlayers, alice.ID)
	assert.Contains(t, g.BlockingPlayers, bob.ID)
	// the players hand should all be developed, and the crib too
	assert.Len(t, g.Hands[alice.ID], 4)
	assert.Contains(t, g.Hands[alice.ID], model.NewCardFromString(`1s`))
	assert.Contains(t, g.Hands[alice.ID], model.NewCardFromString(`2s`))
	assert.Contains(t, g.Hands[alice.ID], model.NewCardFromString(`3s`))
	assert.Contains(t, g.Hands[alice.ID], model.NewCardFromString(`4s`))
	assert.Len(t, g.Hands[bob.ID], 4)
	assert.Contains(t, g.Hands[bob.ID], model.NewCardFromString(`1c`))
	assert.Contains(t, g.Hands[bob.ID], model.NewCardFromString(`2c`))
	assert.Contains(t, g.Hands[bob.ID], model.NewCardFromString(`3c`))
	assert.Contains(t, g.Hands[bob.ID], model.NewCardFromString(`4c`))
	assert.Len(t, g.Crib, 4)
	assert.Contains(t, g.Crib, model.NewCardFromString(`6s`))
	assert.Contains(t, g.Crib, model.NewCardFromString(`6c`))
	assert.Contains(t, g.Crib, model.NewCardFromString(`6d`))
	assert.Contains(t, g.Crib, model.NewCardFromString(`6h`))
	// verify that entering the cutting phase clears out the cut card until it _is_ cut
	assert.NotEqual(t, model.Card{}, g.CutCard)

	aliceAPI.AssertExpectations(t)
	bobAPI.AssertExpectations(t)
}

func TestHandleAction_Pegging(t *testing.T) {
	alice, bob, aliceAPI, bobAPI, abAPIs := testutils.AliceAndBob()

	g := model.Game{
		ID:              model.GameID(5),
		Players:         []model.Player{alice, bob},
		BlockingPlayers: map[model.PlayerID]model.Blocker{bob.ID: model.PegCard},
		CurrentDealer:   alice.ID,
		PlayerColors:    map[model.PlayerID]model.PlayerColor{alice.ID: model.Blue, bob.ID: model.Red},
		CurrentScores:   map[model.PlayerColor]int{model.Blue: 0, model.Red: 0},
		LagScores:       map[model.PlayerColor]int{model.Blue: 0, model.Red: 0},
		Phase:           model.Pegging,
		// cards are chosen so we can test a 31 and a go
		Hands: map[model.PlayerID][]model.Card{
			alice.ID: {
				model.NewCardFromString(`7s`),
				model.NewCardFromString(`8s`),
				model.NewCardFromString(`10s`),
				model.NewCardFromString(`js`),
			},
			bob.ID: {
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
	}

	action := model.PlayerAction{
		GameID:    g.ID,
		ID:        bob.ID,
		Overcomes: model.PegCard,
		Action: model.PegAction{
			Card: g.Hands[bob.ID][0],
		},
	}
	aliceAPI.On(`NotifyBlocking`, model.PegCard, mock.AnythingOfType(`model.Game`), ``).Return(nil).Once()
	err := HandleAction(&g, action, abAPIs)
	assert.Nil(t, err)
	assert.Len(t, g.PeggedCards, 1)
	assert.Contains(t, g.PeggedCards, model.NewPeggedCardFromString(bob.ID, `7c`, g.NumActions()))
	assert.Equal(t, g.CurrentPeg(), 7)

	action = model.PlayerAction{
		GameID:    g.ID,
		ID:        alice.ID,
		Overcomes: model.PegCard,
		Action: model.PegAction{
			Card: g.Hands[alice.ID][0],
		},
	}
	// alice and bob are going to get notified because alice scores a 31 and a run of 3
	aliceAPI.On(`NotifyScoreUpdate`, mock.AnythingOfType(`model.Game`), []string{`pegging`}).Return(nil).Once()
	bobAPI.On(`NotifyScoreUpdate`, mock.AnythingOfType(`model.Game`), []string{`pegging`}).Return(nil).Once()
	bobAPI.On(`NotifyBlocking`, model.PegCard, mock.AnythingOfType(`model.Game`), ``).Return(nil).Once()
	err = HandleAction(&g, action, abAPIs)
	assert.Nil(t, err)
	assert.Len(t, g.PeggedCards, 2)
	assert.Contains(t, g.PeggedCards, model.NewPeggedCardFromString(alice.ID, `7s`, g.NumActions()))
	assert.Equal(t, g.CurrentPeg(), 14)

	action = model.PlayerAction{
		GameID:    g.ID,
		ID:        bob.ID,
		Overcomes: model.PegCard,
		Action: model.PegAction{
			Card: g.Hands[bob.ID][1],
		},
	}

	aliceAPI.On(`NotifyBlocking`, model.PegCard, mock.AnythingOfType(`model.Game`), ``).Return(nil).Once()
	err = HandleAction(&g, action, abAPIs)
	assert.Nil(t, err)
	assert.Len(t, g.PeggedCards, 3)
	assert.Contains(t, g.PeggedCards, model.NewPeggedCardFromString(bob.ID, `9c`, g.NumActions()))
	assert.Equal(t, g.CurrentPeg(), 23)

	// alice tries to peg a card she has already pegged
	action = model.PlayerAction{
		GameID:    g.ID,
		ID:        alice.ID,
		Overcomes: model.PegCard,
		Action: model.PegAction{
			Card: g.Hands[alice.ID][0],
		},
	}
	aliceAPI.On(`NotifyBlocking`, model.PegCard, mock.AnythingOfType(`model.Game`), `Cannot peg same card twice`).Return(nil).Once()
	err = HandleAction(&g, action, abAPIs)
	assert.Nil(t, err)
	assert.Len(t, g.PeggedCards, 3)

	// Alice pegs her 8S and hits 31 with a run of 3
	action = model.PlayerAction{
		GameID:    g.ID,
		ID:        alice.ID,
		Overcomes: model.PegCard,
		Action: model.PegAction{
			Card: g.Hands[alice.ID][1],
		},
	}
	// alice and bob are going to get notified because alice scores a 31 and a run of 3
	aliceAPI.On(`NotifyScoreUpdate`, mock.AnythingOfType(`model.Game`), []string{`pegging`}).Return(nil).Once()
	bobAPI.On(`NotifyScoreUpdate`, mock.AnythingOfType(`model.Game`), []string{`pegging`}).Return(nil).Once()
	bobAPI.On(`NotifyBlocking`, model.PegCard, mock.AnythingOfType(`model.Game`), ``).Return(nil).Once()
	err = HandleAction(&g, action, abAPIs)
	assert.Nil(t, err)
	assert.Len(t, g.PeggedCards, 4)
	assert.Contains(t, g.PeggedCards, model.NewPeggedCardFromString(alice.ID, `8s`, g.NumActions()))
	assert.Equal(t, g.CurrentPeg(), 0)
	assert.Equal(t, g.CurrentScores[g.PlayerColors[alice.ID]], 7)

	// bob pegs his 10C
	action = model.PlayerAction{
		GameID:    g.ID,
		ID:        bob.ID,
		Overcomes: model.PegCard,
		Action: model.PegAction{
			Card: g.Hands[bob.ID][2],
		},
	}
	aliceAPI.On(`NotifyBlocking`, model.PegCard, mock.AnythingOfType(`model.Game`), ``).Return(nil).Once()
	err = HandleAction(&g, action, abAPIs)
	assert.Nil(t, err)
	assert.Len(t, g.PeggedCards, 5)
	assert.Contains(t, g.PeggedCards, model.NewPeggedCardFromString(bob.ID, `10c`, g.NumActions()))
	assert.Equal(t, g.CurrentPeg(), 10)

	action = model.PlayerAction{
		GameID:    g.ID,
		ID:        alice.ID,
		Overcomes: model.PegCard,
		Action: model.PegAction{
			Card: g.Hands[alice.ID][2],
		},
	}
	// alice and bob are going to get notified because alice scores a pair
	aliceAPI.On(`NotifyScoreUpdate`, mock.AnythingOfType(`model.Game`), []string{`pegging`}).Return(nil).Once()
	bobAPI.On(`NotifyScoreUpdate`, mock.AnythingOfType(`model.Game`), []string{`pegging`}).Return(nil).Once()
	bobAPI.On(`NotifyBlocking`, model.PegCard, mock.AnythingOfType(`model.Game`), ``).Return(nil).Once()
	err = HandleAction(&g, action, abAPIs)
	assert.Nil(t, err)
	assert.Len(t, g.PeggedCards, 6)
	assert.Contains(t, g.PeggedCards, model.NewPeggedCardFromString(alice.ID, `10s`, g.NumActions()))
	assert.Equal(t, g.CurrentPeg(), 20)
	assert.Equal(t, g.CurrentScores[g.PlayerColors[alice.ID]], 9)

	action = model.PlayerAction{
		GameID:    g.ID,
		ID:        bob.ID,
		Overcomes: model.PegCard,
		Action: model.PegAction{
			Card: g.Hands[bob.ID][3],
		},
	}

	aliceAPI.On(`NotifyBlocking`, model.PegCard, mock.AnythingOfType(`model.Game`), ``).Return(nil).Once()
	err = HandleAction(&g, action, abAPIs)
	assert.Nil(t, err)
	assert.Len(t, g.PeggedCards, 7)
	assert.Contains(t, g.PeggedCards, model.NewPeggedCardFromString(bob.ID, `jc`, g.NumActions()))
	assert.Equal(t, g.CurrentPeg(), 30)

	// alice tries to peg a jack but she needs to say go
	action = model.PlayerAction{
		GameID:    g.ID,
		ID:        alice.ID,
		Overcomes: model.PegCard,
		Action: model.PegAction{
			Card: g.Hands[alice.ID][3],
		},
	}
	aliceAPI.On(`NotifyBlocking`, model.PegCard, mock.AnythingOfType(`model.Game`), `Cannot peg card with this value`).Return(nil).Once()
	err = HandleAction(&g, action, abAPIs)
	assert.Nil(t, err)
	assert.Len(t, g.PeggedCards, 7)

	// alice tries to peg a card she doesn't have
	action = model.PlayerAction{
		GameID:    g.ID,
		ID:        alice.ID,
		Overcomes: model.PegCard,
		Action: model.PegAction{
			Card: model.NewCardFromString(`2h`),
		},
	}
	aliceAPI.On(`NotifyBlocking`, model.PegCard, mock.AnythingOfType(`model.Game`), `Cannot peg card you don't have`).Return(nil).Once()
	err = HandleAction(&g, action, abAPIs)
	assert.Nil(t, err)
	assert.Len(t, g.PeggedCards, 7)

	// alice says go
	action = model.PlayerAction{
		GameID:    g.ID,
		ID:        alice.ID,
		Overcomes: model.PegCard,
		Action: model.PegAction{
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

	// Current peg has reset to 0
	assert.Equal(t, 0, g.CurrentPeg())

	// alice scores last card
	action = model.PlayerAction{
		GameID:    g.ID,
		ID:        alice.ID,
		Overcomes: model.PegCard,
		Action: model.PegAction{
			Card: g.Hands[alice.ID][3],
		},
	}
	aliceAPI.On(`NotifyScoreUpdate`, mock.AnythingOfType(`model.Game`), []string{`last card`}).Return(nil).Once()
	bobAPI.On(`NotifyScoreUpdate`, mock.AnythingOfType(`model.Game`), []string{`last card`}).Return(nil).Once()
	// bob will be up to count his hand
	bobAPI.On(`NotifyBlocking`, model.CountHand, mock.AnythingOfType(`model.Game`), ``).Return(nil).Once()
	err = HandleAction(&g, action, abAPIs)
	assert.Nil(t, err)
	assert.Len(t, g.PeggedCards, 8)
	assert.Contains(t, g.PeggedCards, model.NewPeggedCardFromString(alice.ID, `js`, g.NumActions()))
	assert.Equal(t, 0, g.CurrentPeg(), `should not have a "current peg" in another phase`)
	assert.Equal(t, 10, g.CurrentScores[g.PlayerColors[alice.ID]])

	// we have moved on to counting hands, and bob is up
	assert.Equal(t, model.Counting, g.Phase)
	assert.Contains(t, g.BlockingPlayers, bob.ID)
	assert.NotContains(t, g.BlockingPlayers, alice.ID)

	aliceAPI.AssertExpectations(t)
	bobAPI.AssertExpectations(t)
}

func TestHandleAction_Counting(t *testing.T) {
	alice, bob, aliceAPI, bobAPI, abAPIs := testutils.AliceAndBob()

	g := model.Game{
		ID:              model.GameID(5),
		Players:         []model.Player{alice, bob},
		BlockingPlayers: map[model.PlayerID]model.Blocker{bob.ID: model.CountHand},
		CurrentDealer:   alice.ID,
		PlayerColors:    map[model.PlayerID]model.PlayerColor{alice.ID: model.Blue, bob.ID: model.Red},
		CurrentScores:   map[model.PlayerColor]int{model.Blue: 0, model.Red: 0},
		LagScores:       map[model.PlayerColor]int{model.Blue: 0, model.Red: 0},
		Phase:           model.Counting,
		Hands: map[model.PlayerID][]model.Card{
			alice.ID: {
				model.NewCardFromString(`7s`),
				model.NewCardFromString(`8s`),
				model.NewCardFromString(`9s`),
				model.NewCardFromString(`10s`),
			},
			bob.ID: {
				model.NewCardFromString(`7c`),
				model.NewCardFromString(`8c`),
				model.NewCardFromString(`9c`),
				model.NewCardFromString(`10c`),
			},
		},
		CutCard:     model.NewCardFromString(`7h`),
		Crib:        make([]model.Card, 4),
		PeggedCards: make([]model.PeggedCard, 0, 8),
	}

	action := model.PlayerAction{
		GameID:    g.ID,
		ID:        bob.ID,
		Overcomes: model.CountHand,
		Action: model.CountHandAction{
			Pts: 100,
		},
	}
	bobAPI.On(`NotifyBlocking`, model.CountHand, mock.AnythingOfType(`model.Game`), `you did not submit the correct number of points for your hand`).Return(nil).Once()
	err := HandleAction(&g, action, abAPIs)
	assert.Error(t, err)
	assert.EqualError(t, err, `wrong number of points`)
	assert.Equal(t, 0, g.CurrentScores[g.PlayerColors[bob.ID]])
	assert.NotContains(t, g.BlockingPlayers, alice.ID)
	assert.Contains(t, g.BlockingPlayers, bob.ID)

	action = model.PlayerAction{
		GameID:    g.ID,
		ID:        bob.ID,
		Overcomes: model.CountHand,
		Action: model.CountHandAction{
			Pts: 18,
		},
	}
	bobAPI.On(`NotifyScoreUpdate`, mock.AnythingOfType(`model.Game`), []string{`hand (7H: 7C, 8C, 9C, 10C)`}).Return(nil).Once()
	aliceAPI.On(`NotifyScoreUpdate`, mock.AnythingOfType(`model.Game`), []string{`hand (7H: 7C, 8C, 9C, 10C)`}).Return(nil).Once()
	aliceAPI.On(`NotifyBlocking`, model.CountHand, mock.AnythingOfType(`model.Game`), ``).Return(nil).Once()
	err = HandleAction(&g, action, abAPIs)
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
	bobAPI.On(`NotifyScoreUpdate`, mock.AnythingOfType(`model.Game`), []string{`hand (7H: 7S, 8S, 9S, 10S)`}).Return(nil).Once()
	aliceAPI.On(`NotifyScoreUpdate`, mock.AnythingOfType(`model.Game`), []string{`hand (7H: 7S, 8S, 9S, 10S)`}).Return(nil).Once()
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

func TestHandleAction_CribCounting(t *testing.T) {
	alice, bob, aliceAPI, bobAPI, abAPIs := testutils.AliceAndBob()

	g := model.Game{
		ID:              model.GameID(5),
		Players:         []model.Player{alice, bob},
		BlockingPlayers: map[model.PlayerID]model.Blocker{alice.ID: model.CountCrib},
		CurrentDealer:   alice.ID,
		PlayerColors:    map[model.PlayerID]model.PlayerColor{alice.ID: model.Blue, bob.ID: model.Red},
		CurrentScores:   map[model.PlayerColor]int{model.Blue: 0, model.Red: 0},
		LagScores:       map[model.PlayerColor]int{model.Blue: 0, model.Red: 0},
		Phase:           model.CribCounting,
		Hands:           make(map[model.PlayerID][]model.Card, 2),
		CutCard:         model.NewCardFromString(`7h`),
		Crib: []model.Card{
			model.NewCardFromString(`7s`),
			model.NewCardFromString(`8s`),
			model.NewCardFromString(`9s`),
			model.NewCardFromString(`10s`),
		},
		PeggedCards: make([]model.PeggedCard, 0, 8),
	}

	action := model.PlayerAction{
		GameID:    g.ID,
		ID:        alice.ID,
		Overcomes: model.CountCrib,
		Action: model.CountCribAction{
			Pts: 100,
		},
	}
	aliceAPI.On(`NotifyBlocking`, model.CountCrib, mock.AnythingOfType(`model.Game`), `you did not submit the correct number of points for the crib`).Return(nil).Once()
	err := HandleAction(&g, action, abAPIs)
	assert.Error(t, err)
	assert.EqualError(t, err, `wrong number of points`)
	assert.Equal(t, 0, g.CurrentScores[g.PlayerColors[alice.ID]])
	assert.NotContains(t, g.BlockingPlayers, bob.ID)
	assert.Contains(t, g.BlockingPlayers, alice.ID)

	action = model.PlayerAction{
		GameID:    g.ID,
		ID:        alice.ID,
		Overcomes: model.CountCrib,
		Action: model.CountCribAction{
			Pts: 14,
		},
	}
	bobAPI.On(`NotifyScoreUpdate`, mock.AnythingOfType(`model.Game`), []string{`crib (7H: 7S, 8S, 9S, 10S)`}).Return(nil).Once()
	aliceAPI.On(`NotifyScoreUpdate`, mock.AnythingOfType(`model.Game`), []string{`crib (7H: 7S, 8S, 9S, 10S)`}).Return(nil).Once()
	bobAPI.On(`NotifyBlocking`, model.DealCards, mock.AnythingOfType(`model.Game`), ``).Return(nil).Once()
	err = HandleAction(&g, action, abAPIs)
	assert.Nil(t, err)
	assert.Equal(t, 14, g.CurrentScores[g.PlayerColors[alice.ID]])
	assert.Contains(t, g.BlockingPlayers, bob.ID)
	assert.NotContains(t, g.BlockingPlayers, alice.ID)

	aliceAPI.AssertExpectations(t)
	bobAPI.AssertExpectations(t)
}

func TestHandleAction_DealAgain(t *testing.T) {
	alice, bob, aliceAPI, bobAPI, abAPIs := testutils.AliceAndBob()

	// Start handlers get called in the *Ready phases, so start from crib
	// counting to make sure we pass through the DealReady phase
	g := model.Game{
		ID:              model.GameID(5),
		Players:         []model.Player{alice, bob},
		BlockingPlayers: map[model.PlayerID]model.Blocker{alice.ID: model.CountCrib},
		CurrentDealer:   alice.ID,
		PlayerColors:    map[model.PlayerID]model.PlayerColor{alice.ID: model.Blue, bob.ID: model.Red},
		CurrentScores:   map[model.PlayerColor]int{model.Blue: 0, model.Red: 0},
		LagScores:       map[model.PlayerColor]int{model.Blue: 0, model.Red: 0},
		Phase:           model.CribCounting,
		Hands: map[model.PlayerID][]model.Card{
			alice.ID: {
				model.NewCardFromString(`7s`),
				model.NewCardFromString(`8s`),
				model.NewCardFromString(`9s`),
				model.NewCardFromString(`10s`),
			},
			bob.ID: {
				model.NewCardFromString(`7c`),
				model.NewCardFromString(`8c`),
				model.NewCardFromString(`9c`),
				model.NewCardFromString(`10c`),
			},
		},
		CutCard: model.NewCardFromString(`7h`),
		Crib: []model.Card{
			model.NewCardFromString(`7d`),
			model.NewCardFromString(`8d`),
			model.NewCardFromString(`9d`),
			model.NewCardFromString(`10d`),
		},
		PeggedCards: make([]model.PeggedCard, 0, 8),
	}

	action := model.PlayerAction{
		GameID:    g.ID,
		ID:        alice.ID,
		Overcomes: model.CountCrib,
		Action: model.CountCribAction{
			Pts: 14,
		},
	}
	bobAPI.On(`NotifyScoreUpdate`, mock.AnythingOfType(`model.Game`), []string{`crib (7H: 7D, 8D, 9D, 10D)`}).Return(nil).Once()
	aliceAPI.On(`NotifyScoreUpdate`, mock.AnythingOfType(`model.Game`), []string{`crib (7H: 7D, 8D, 9D, 10D)`}).Return(nil).Once()
	bobAPI.On(`NotifyBlocking`, model.DealCards, mock.AnythingOfType(`model.Game`), ``).Return(nil).Once()
	err := HandleAction(&g, action, abAPIs)
	require.Nil(t, err)

	action = model.PlayerAction{
		GameID:    g.ID,
		ID:        bob.ID,
		Overcomes: model.DealCards,
		Action: model.DealAction{
			NumShuffles: 1,
		},
	}
	aliceAPI.On(`NotifyMessage`, mock.AnythingOfType(`model.Game`), mock.MatchedBy(func(s string) bool { return strings.HasPrefix(s, `Received Hand `) })).Return(nil).Once()
	aliceAPI.On(`NotifyBlocking`, model.CribCard, mock.AnythingOfType(`model.Game`), `needs to cut 2 cards`).Return(nil).Once()
	bobAPI.On(`NotifyMessage`, mock.AnythingOfType(`model.Game`), mock.MatchedBy(func(s string) bool { return strings.HasPrefix(s, `Received Hand `) })).Return(nil).Once()
	bobAPI.On(`NotifyBlocking`, model.CribCard, mock.AnythingOfType(`model.Game`), `needs to cut 2 cards`).Return(nil).Once()

	err = HandleAction(&g, action, abAPIs)
	assert.Nil(t, err)
	assert.Len(t, g.Hands[alice.ID], 6)
	assert.Len(t, g.Hands[bob.ID], 6)

	aliceAPI.AssertExpectations(t)
	bobAPI.AssertExpectations(t)
}
