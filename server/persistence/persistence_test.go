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
	"github.com/joshprzybyszewski/cribbage/server/persistence/mongodb"
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
		// We assume you have mongodb stood up locally when running without -short
		mongo, err := mongodb.New(context.Background(), ``)
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
	alice, bob, _ := testutils.EmptyAliceAndBob()

	g1 := model.Game{
		ID:              model.NewGameID(),
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

	require.NoError(t, db.SaveGame(g1))

	actGame, err := db.GetGame(g1.ID)
	require.NoError(t, err)
	assert.Equal(t, g1Copy, actGame)
}

func testSaveGameMultipleTimes(t *testing.T, db persistence.DB) {
	alice, bob, abAPIs := testutils.EmptyAliceAndBob()

	checkPersistedGame := func(expGame model.Game) {
		actGame, err := db.GetGame(expGame.ID)
		require.NoError(t, err)
		assert.Equal(t, expGame, actGame)
	}

	g, err := play.CreateGame([]model.Player{alice, bob}, abAPIs)
	require.NoError(t, err)

	_, err = db.GetGame(g.ID)
	require.Error(t, err)
	assert.EqualError(t, err, persistence.ErrGameNotFound.Error())

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
			Info: int32(8383),
		}},
	}
	p1Copy := p1

	assert.NoError(t, db.SaveInteraction(p1))

	actPM, err := db.GetInteraction(p1.PlayerID)
	require.NoError(t, err)
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

	// Right now, CreateGame assigns colors to players, but we may
	// not always do that. Clear out that map but save it so that we
	// _can_ assign colors:)
	playerColors := g.PlayerColors
	g.PlayerColors = nil

	require.NoError(t, db.SaveGame(g))

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
		require.NoError(t, err)
		assert.Equal(t, expGame, actGame)
	}

	g, err := play.CreateGame([]model.Player{alice, bob}, abAPIs)
	require.NoError(t, err)

	_, err = db.GetGame(g.ID)
	require.Error(t, err)
	assert.EqualError(t, err, persistence.ErrGameNotFound.Error())

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

func TestTransactionality(t *testing.T) {
	// Right now (and probably ever) the memory persistence isn't transactional
	dbs := map[string]func() persistence.DB{
		// `memory`: func() persistence.DB { return inmem },
	}

	if !testing.Short() {
		dbs[`mongodb`] = func() persistence.DB {
			// We assume you have mongodb stood up locally when running without -short
			mongo, err := mongodb.New(context.Background(), ``)
			require.NoError(t, err)
			return mongo
		}
	}

	txTests := map[string]txTest{
		`player service`:              playerTxTest,
		`game service`:                gameTxTest,
		`rollback the player service`: rollbackPlayerTxTest,
	}

	for dbName, db := range dbs {
		for testName, txTest := range txTests {
			t.Run(dbName+`:`+testName, func(t1 *testing.T) { txTest(t1, db(), db(), db()) })
		}
	}
}

type txTest func(t *testing.T, db1, db2, postCommitDB persistence.DB)

func playerTxTest(t *testing.T, db1, db2, postCommitDB persistence.DB) {
	require.NoError(t, db1.Start())
	require.NoError(t, db2.Start())

	p1 := model.Player{
		ID:    model.PlayerID(rand.String(50)),
		Name:  `player 1`,
		Games: map[model.GameID]model.PlayerColor{},
	}

	assert.NoError(t, db1.CreatePlayer(p1))
	p1Mod := p1
	p1Mod.Name = `different player 1 name`
	assert.NoError(t, db2.CreatePlayer(p1Mod))

	err := db1.CreatePlayer(p1)
	assert.EqualError(t, err, persistence.ErrPlayerAlreadyExists.Error())

	savedP1, err := db1.GetPlayer(p1.ID)
	require.NoError(t, err)
	assert.Equal(t, p1, savedP1)

	savedP1Mod, err := db2.GetPlayer(p1.ID)
	require.NoError(t, err)
	assert.NotEqual(t, p1, savedP1Mod)

	assert.NoError(t, db1.Commit())
	// the second connection tried to save a different player one, so committing should error
	// assert.Error(t, db2.Commit())
	// assert.NoError(t, db2.Rollback())

	postCommitP1, err := postCommitDB.GetPlayer(p1.ID)
	require.NoError(t, err)
	assert.Equal(t, p1, postCommitP1)
	assert.NotEqual(t, p1Mod, postCommitP1)
}

func rollbackPlayerTxTest(t *testing.T, db1, db2, postCommitDB persistence.DB) {
	require.NoError(t, db1.Start())
	require.NoError(t, db2.Start())

	p1 := model.Player{
		ID:    model.PlayerID(rand.String(50)),
		Name:  `player 1`,
		Games: map[model.GameID]model.PlayerColor{},
	}

	assert.NoError(t, db1.CreatePlayer(p1))
	p2 := model.Player{
		ID:    model.PlayerID(rand.String(50)),
		Name:  `player 2`,
		Games: map[model.GameID]model.PlayerColor{},
	}
	assert.NoError(t, db2.CreatePlayer(p2))

	savedP1, err := db1.GetPlayer(p1.ID)
	require.NoError(t, err)
	assert.Equal(t, p1, savedP1)

	savedP2, err := db1.GetPlayer(p2.ID)
	assert.Error(t, err)
	assert.NotEqual(t, p2, savedP2)

	savedP1, err = db2.GetPlayer(p1.ID)
	assert.Error(t, err)
	assert.NotEqual(t, p1, savedP1)

	savedP2, err = db2.GetPlayer(p2.ID)
	require.NoError(t, err)
	assert.Equal(t, p2, savedP2)

	assert.NoError(t, db1.Rollback())
	assert.NoError(t, db2.Rollback())

	postCommitP1, err := postCommitDB.GetPlayer(p1.ID)
	assert.Error(t, err)
	assert.NotEqual(t, p1, postCommitP1)

	postCommitP2, err := postCommitDB.GetPlayer(p2.ID)
	assert.Error(t, err)
	assert.NotEqual(t, p2, postCommitP2)

}

func gameTxTest(t *testing.T, db1, db2, postCommitDB persistence.DB) {
	require.NoError(t, db1.Start())
	require.NoError(t, db2.Start())

	alice, bob, _ := testutils.EmptyAliceAndBob()

	g1 := model.Game{
		ID:              model.NewGameID(),
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

	require.NoError(t, db1.SaveGame(g1))

	actGame, err := db1.GetGame(g1.ID)
	require.NoError(t, err)
	assert.Equal(t, g1Copy, actGame)

	actGame, err = db2.GetGame(g1.ID)
	assert.Error(t, err)
	assert.NotEqual(t, g1Copy, actGame)

	assert.NoError(t, db1.Commit())
	// the second connection tried to save a different game one, so committing should error
	// assert.Error(t, db2.Commit())
	// assert.NoError(t, db2.Rollback())

	postCommitGame, err := postCommitDB.GetGame(g1.ID)
	require.NoError(t, err)
	assert.Equal(t, g1Copy, postCommitGame)
}
