package persistence_test

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
	"github.com/joshprzybyszewski/cribbage/server/persistence/memory"
	"github.com/joshprzybyszewski/cribbage/server/persistence/mongodb"
	"github.com/joshprzybyszewski/cribbage/server/play"
	"github.com/joshprzybyszewski/cribbage/utils/rand"
	"github.com/joshprzybyszewski/cribbage/utils/testutils"
)

type dbTest func(*testing.T, persistence.DB)

var (
	tests = map[string]dbTest{
		`createPlayer`:    testCreatePlayer,
		`saveGame`:        testSaveGame,
		`resaveGame`:      testSaveGameMultipleTimes,
		`saveInteraction`: testSaveInteraction,
	}
)

func setup() (a, b model.Player, am, bm *interaction.Mock, pAPIs map[model.PlayerID]interaction.Player) {
	alice, bob, aAPI, bAPI, abAPIs := testutils.AliceAndBob()

	aAPI.On(`ID`).Return(alice.ID)
	aAPI.On(`NotifyBlocking`, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	aAPI.On(`NotifyMessage`, mock.Anything, mock.Anything).Return(nil)
	aAPI.On(`NotifyScoreUpdate`, mock.Anything, mock.Anything).Return(nil)

	bAPI.On(`ID`).Return(bob.ID)
	bAPI.On(`NotifyBlocking`, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	bAPI.On(`NotifyMessage`, mock.Anything, mock.Anything).Return(nil)
	bAPI.On(`NotifyScoreUpdate`, mock.Anything, mock.Anything).Return(nil)
	return alice, bob, aAPI, bAPI, abAPIs
}

func persistenceGameCopy(dst *model.Game, src model.Game) {
	*dst = src

	// The deck will always be nil after a copy
	dst.Deck = nil

	dst.Players = make([]model.Player, len(src.Players))
	_ = copy(dst.Players, src.Players)

	dst.BlockingPlayers = make(map[model.PlayerID]model.Blocker, len(src.BlockingPlayers))
	for k, v := range src.BlockingPlayers {
		dst.BlockingPlayers[k] = v
	}

	dst.PlayerColors = make(map[model.PlayerID]model.PlayerColor, len(src.PlayerColors))
	for k, v := range src.PlayerColors {
		dst.PlayerColors[k] = v
	}

	dst.CurrentScores = make(map[model.PlayerColor]int, len(src.CurrentScores))
	for k, v := range src.CurrentScores {
		dst.CurrentScores[k] = v
	}

	dst.LagScores = make(map[model.PlayerColor]int, len(src.LagScores))
	for k, v := range src.LagScores {
		dst.LagScores[k] = v
	}

	dst.Hands = make(map[model.PlayerID][]model.Card, len(src.Hands))
	for k, v := range src.Hands {
		newHand := make([]model.Card, len(v))
		copy(newHand, v)
		dst.Hands[k] = newHand
	}

	dst.Crib = make([]model.Card, len(src.Crib))
	_ = copy(dst.Crib, src.Crib)

	dst.PeggedCards = make([]model.PeggedCard, len(src.PeggedCards))
	_ = copy(dst.PeggedCards, src.PeggedCards)
}

func TestDB(t *testing.T) {
	dbs := map[string]persistence.DB{
		`memory`: memory.New(),
	}

	// TODO get mongodb tests running in travis?
	if !testing.Short() {
		// We assume you have mongodb stood up locally when running without -short
		mongo, err := mongodb.New(``)
		require.NoError(t, err)

		dbs[`mongodb`] = mongo
	}

	for dbName, db := range dbs {
		for testName, testFn := range tests {
			t.Run(dbName+`:`+testName, func(t1 *testing.T) { testFn(t1, db) })
		}
	}
}

func testCreatePlayer(t *testing.T, db persistence.DB) {
	p1 := model.Player{
		ID:    model.PlayerID(rand.String(50)),
		Name:  `player 1`,
		Games: map[model.GameID]model.PlayerColor{},
	}

	assert.NoError(t, db.CreatePlayer(p1))

	err := db.CreatePlayer(p1)
	assert.EqualError(t, err, persistence.ErrPlayerAlreadyExists.Error())

	p1conflict := model.Player{
		ID:   p1.ID,
		Name: `different name`,
		Games: map[model.GameID]model.PlayerColor{
			model.GameID(4): model.Blue,
		},
	}
	err = db.CreatePlayer(p1conflict)
	assert.EqualError(t, err, persistence.ErrPlayerAlreadyExists.Error())

	p2 := model.Player{
		ID:   model.PlayerID(rand.String(50)),
		Name: `player 2`,
		Games: map[model.GameID]model.PlayerColor{
			model.GameID(1825): model.Blue,
			model.GameID(26):   model.Red,
			model.GameID(33):   model.Green,
			model.GameID(108):  model.Red,
		},
	}
	p2Copy := p2
	// Don't keep the same memory space for the games copy
	p2Copy.Games = map[model.GameID]model.PlayerColor{
		model.GameID(1825): model.Blue,
		model.GameID(26):   model.Red,
		model.GameID(33):   model.Green,
		model.GameID(108):  model.Red,
	}

	assert.NoError(t, db.CreatePlayer(p2))

	actP2, err := db.GetPlayer(p2.ID)
	require.NoError(t, err)
	assert.Equal(t, p2Copy, actP2)
}

func testSaveGame(t *testing.T, db persistence.DB) {
	alice, bob, _, _, _ := setup()

	g1 := model.Game{
		ID:              model.GameID(rand.Intn(1000)),
		Players:         []model.Player{alice, bob},
		Deck:            model.NewDeck(),
		BlockingPlayers: map[model.PlayerID]model.Blocker{bob.ID: model.PegCard},
		CurrentDealer:   alice.ID,
		PlayerColors:    map[model.PlayerID]model.PlayerColor{alice.ID: model.Blue, bob.ID: model.Red},
		CurrentScores:   map[model.PlayerColor]int{model.Blue: 0, model.Red: 0},
		LagScores:       map[model.PlayerColor]int{model.Blue: 0, model.Red: 0},
		Phase:           model.Pegging,
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
	g1Copy := g1

	require.NoError(t, db.SaveGame(g1))

	actGame, err := db.GetGame(g1.ID)
	require.NoError(t, err)
	// the deck should be nil
	g1Copy.Deck = nil
	assert.Equal(t, g1Copy, actGame)
}

func testSaveGameMultipleTimes(t *testing.T, db persistence.DB) {
	alice, bob, _, _, abAPIs := setup()

	checkPersistedGame := func(expGame model.Game) {
		actGame, err := db.GetGame(expGame.ID)
		require.NoError(t, err)
		// the deck should always be nil
		expGame.Deck = nil
		assert.Equal(t, expGame, actGame)
	}

	g, err := play.CreateGame([]model.Player{alice, bob}, abAPIs)
	require.NoError(t, err)
	var gCopy model.Game
	persistenceGameCopy(&gCopy, g)

	require.NoError(t, db.SaveGame(g))

	checkPersistedGame(gCopy)

	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        alice.ID,
		GameID:    g.ID,
		Overcomes: model.DealCards,
		Action:    model.DealAction{NumShuffles: 10},
	}, abAPIs))
	persistenceGameCopy(&gCopy, g)

	require.NoError(t, db.SaveGame(g))
	checkPersistedGame(gCopy)

	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        alice.ID,
		GameID:    g.ID,
		Overcomes: model.CribCard,
		Action:    model.BuildCribAction{Cards: []model.Card{g.Hands[alice.ID][0], g.Hands[alice.ID][1]}},
	}, abAPIs))
	persistenceGameCopy(&gCopy, g)

	require.NoError(t, db.SaveGame(g))
	checkPersistedGame(gCopy)

	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        bob.ID,
		GameID:    g.ID,
		Overcomes: model.CribCard,
		Action:    model.BuildCribAction{Cards: []model.Card{g.Hands[bob.ID][0], g.Hands[bob.ID][1]}},
	}, abAPIs))
	persistenceGameCopy(&gCopy, g)

	require.NoError(t, db.SaveGame(g))
	checkPersistedGame(gCopy)
}

func testSaveInteraction(t *testing.T, db persistence.DB) {
	p1 := interaction.PlayerMeans{
		PlayerID:      model.PlayerID(rand.String(50)),
		PreferredMode: interaction.Localhost,
		Interactions: []interaction.Means{{
			Mode: interaction.Localhost,
			Info: 8383,
		}},
	}
	p1Copy := p1

	assert.NoError(t, db.SaveInteraction(p1))

	actPM, err := db.GetInteraction(p1.PlayerID)
	assert.NoError(t, err)
	assert.Equal(t, p1Copy, actPM)

	assert.NoError(t, db.SaveInteraction(p1))

	p1update := interaction.PlayerMeans{
		PlayerID:      p1.PlayerID,
		PreferredMode: interaction.Localhost,
		Interactions: []interaction.Means{{
			Mode: interaction.Localhost,
			Info: 8484,
		}},
	}
	assert.NoError(t, db.SaveInteraction(p1update))

	actPM, err = db.GetInteraction(p1.PlayerID)
	assert.NoError(t, err)
	assert.NotEqual(t, p1Copy, actPM)
}
