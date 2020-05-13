package persistence_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
	"github.com/joshprzybyszewski/cribbage/server/persistence/memory"
	"github.com/joshprzybyszewski/cribbage/server/persistence/mysql"
	"github.com/joshprzybyszewski/cribbage/server/play"
	"github.com/joshprzybyszewski/cribbage/utils/rand"
	"github.com/joshprzybyszewski/cribbage/utils/testutils"
)

type dbTest func(*testing.T, persistence.DB)

var (
	tests = map[string]dbTest{
		`createPlayer`:          testCreatePlayer,
		`saveGame`:              testSaveGame,
		`resaveGame`:            testSaveGameMultipleTimes,
		`saveGameMissingAction`: testSaveGameWithMissingAction,
		`saveInteraction`:       testSaveInteraction,
		`addColorToGame`:        testAddPlayerColorToGame,
	}
)

func persistenceGameCopy(dst *model.Game, src model.Game) {
	*dst = src

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

	if !testing.Short() {
		/* Skip mongodb tests for now
		// We assume you have mongodb stood up locally when running without -short
		mongo, err := mongodb.New(context.Background(), ``)
		require.NoError(t, err)

		dbs[`mongodb`] = mongo
		*/

		// We further assume you have mysql stood up locally when running without -short
		cfg := mysql.Config{
			DSNUser:      `root`, // travis ci uses either the root or the travis user
			DSNPassword:  ``,
			DSNHost:      `127.0.0.1`,
			DSNPort:      3306,
			DatabaseName: `testing_cribbage`,
			DSNParams:    ``,
		}
		mysqlDB, err := mysql.New(context.Background(), cfg)
		require.NoError(t, err)

		dbs[`mysql`] = mysqlDB
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
		ID:    model.PlayerID(rand.String(50)),
		Name:  `player 2`,
		Games: map[model.GameID]model.PlayerColor{},
	}
	p2Copy := p2

	assert.NoError(t, db.CreatePlayer(p2))

	actP2, err := db.GetPlayer(p2.ID)
	require.NoError(t, err)
	assert.Equal(t, p2Copy, actP2)
}

func testSaveGame(t *testing.T, db persistence.DB) {
	alice, bob, _ := testutils.EmptyAliceAndBob()

	g1 := model.Game{
		ID:              model.GameID(rand.Intn(1000)),
		Players:         []model.Player{alice, bob},
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

	for _, p := range g1.Players {
		require.NoError(t, db.CreatePlayer(p))
	}
	require.NoError(t, db.CreateGame(g1))

	actGame, err := db.GetGame(g1.ID)
	require.NoError(t, err, `expected to find game with id "%d"`, g1.ID)
	assert.Equal(t, g1Copy, actGame)
}

func testSaveGameMultipleTimes(t *testing.T, db persistence.DB) {
	alice, bob, abAPIs := testutils.EmptyAliceAndBob()

	checkPersistedGame := func(expGame model.Game) {
		actGame, err := db.GetGame(expGame.ID)
		require.NoError(t, err, `expected to find game with id "%d"`, expGame.ID)
		assert.Equal(t, expGame, actGame)
	}

	g, err := play.CreateGame([]model.Player{alice, bob}, abAPIs)
	require.NoError(t, err)

	for _, p := range g.Players {
		require.NoError(t, db.CreatePlayer(p))
	}

	_, err = db.GetGame(g.ID)
	require.Error(t, err)
	assert.EqualError(t, err, persistence.ErrGameNotFound.Error())

	var gCopy model.Game
	persistenceGameCopy(&gCopy, g)

	require.NoError(t, db.CreateGame(g))

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
			Info: string(`8383`),
		}},
	}
	p1Copy := p1

	assert.NoError(t, db.SaveInteraction(p1))

	actPM, err := db.GetInteraction(p1.PlayerID)
	require.NoError(t, err, `expected to find player with id "%s"`, p1.PlayerID)
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

func testAddPlayerColorToGame(t *testing.T, db persistence.DB) {
	alice, bob, abAPIs := testutils.EmptyAliceAndBob()
	require.NoError(t, db.CreatePlayer(alice))
	require.NoError(t, db.CreatePlayer(bob))

	g, err := play.CreateGame([]model.Player{alice, bob}, abAPIs)
	require.NoError(t, err)
	for _, p := range g.Players {
		require.NoError(t, db.CreatePlayer(p))
	}

	// Right now, CreateGame assigns colors to players, but we may
	// not always do that. Clear out that map but save it so that we
	// _can_ assign colors:)
	playerColors := g.PlayerColors
	g.PlayerColors = nil

	require.NoError(t, db.CreateGame(g))

	for _, pID := range []model.PlayerID{alice.ID, bob.ID} {
		require.NoError(t, db.AddPlayerColorToGame(pID, playerColors[pID], g.ID))
	}

	a2, err := db.GetPlayer(alice.ID)
	require.NoError(t, err)
	assert.NotEqual(t, alice, a2)
	assert.Equal(t, playerColors[alice.ID], a2.Games[g.ID])

	b2, err := db.GetPlayer(bob.ID)
	require.NoError(t, err)
	assert.NotEqual(t, bob, b2)
	assert.Equal(t, playerColors[bob.ID], b2.Games[g.ID])

	g2, err := db.GetGame(g.ID)
	require.NoError(t, err)
	assert.NotEqual(t, g, g2)
	assert.Equal(t, g2.PlayerColors[alice.ID], a2.Games[g.ID])
	assert.Equal(t, g2.PlayerColors[bob.ID], b2.Games[g.ID])
}

func testSaveGameWithMissingAction(t *testing.T, db persistence.DB) {
	alice, bob, abAPIs := testutils.EmptyAliceAndBob()

	checkPersistedGame := func(expGame model.Game) {
		actGame, err := db.GetGame(expGame.ID)
		require.NoError(t, err, `expected to find game with id "%d"`, expGame.ID)
		assert.Equal(t, expGame, actGame)
	}

	g, err := play.CreateGame([]model.Player{alice, bob}, abAPIs)
	require.NoError(t, err)
	for _, p := range g.Players {
		require.NoError(t, db.CreatePlayer(p))
	}

	_, err = db.GetGame(g.ID)
	require.Error(t, err)
	assert.EqualError(t, err, persistence.ErrGameNotFound.Error())

	var gCopy model.Game
	persistenceGameCopy(&gCopy, g)

	require.NoError(t, db.CreateGame(g))

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

	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        bob.ID,
		GameID:    g.ID,
		Overcomes: model.CutCard,
		Action:    model.CutDeckAction{Percentage: 0.600},
	}, abAPIs))
	persistenceGameCopy(&gCopy, g)

	// corrupt an action
	badAction := g.Actions[1]
	badAction.Action = model.CountCribAction{Pts: 100}
	badAction.Overcomes = model.CountCrib
	g.Actions[1] = badAction
	require.Error(t, db.SaveGame(g), `saving a game with a corrupted action is a :badtime:`)

	// splice out an action
	g.Actions = append(g.Actions[:1], g.Actions[2:]...)
	require.Error(t, db.SaveGame(g), `saving a game with a missing action is a :badtime:`)
}
