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
		Crib:            make([]model.Card, 0, 4),
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
	// now the game is blocked by both players needing to build the crib
	require.Len(t, g.BlockingPlayers, 2)
	assert.Contains(t, g.BlockingPlayers, alice.ID)
	assert.Contains(t, g.BlockingPlayers, bob.ID)
	// the players should have 6 card hands
	assert.Len(t, g.Hands[alice.ID], 6)
	assert.Len(t, g.Hands[bob.ID], 6)

	aliceAPI.AssertExpectations(t)
	bobAPI.AssertExpectations(t)
}
