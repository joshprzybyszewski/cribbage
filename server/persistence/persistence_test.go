package persistence_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
	"github.com/joshprzybyszewski/cribbage/server/persistence/memory"
	"github.com/joshprzybyszewski/cribbage/server/persistence/mongodb"
	"github.com/joshprzybyszewski/cribbage/server/persistence/mysql"
	"github.com/joshprzybyszewski/cribbage/server/play"
	"github.com/joshprzybyszewski/cribbage/utils/rand"
	"github.com/joshprzybyszewski/cribbage/utils/testutils"
)

type dbName string
type dbTest func(*testing.T, dbName, persistence.DB)

const (
	memoryDB dbName = `memoryDB`
	mongoDB  dbName = `mongoDB`
	mysqlDB  dbName = `mysqlDB`
)

var (
	tests = map[string]dbTest{
		`createPlayer`:                  testCreatePlayer,
		`createPlayersWithSimilarNames`: testCreatePlayersWithSimilarNames,
		`saveGame`:                      testCreateGame,
		`resaveGame`:                    testSaveGameMultipleTimes,
		`saveGameMissingAction`:         testSaveGameWithMissingAction,
		`saveInteraction`:               testSaveInteraction,
		`addColorToGame`:                testAddPlayerColorToGame,
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

	dst.Actions = make([]model.PlayerAction, len(src.Actions))
	_ = copy(dst.Actions, src.Actions)
}

func checkPersistedGame(t *testing.T, name dbName, db persistence.DB, expGame model.Game) {
	actGame, err := db.GetGame(expGame.ID)
	require.NoError(t, err, `expected to find game with id "%d"`, expGame.ID)
	if len(actGame.Crib) == 0 {
		expGame.Crib = nil
		actGame.Crib = nil
	}
	if len(actGame.Actions) == 0 {
		expGame.Actions = nil
		actGame.Actions = nil
	}
	for i := range actGame.Actions {
		if !(name == memoryDB || name == mongoDB) {
			// memory provider and mongodb do not have this feature implemented
			assert.NotEqual(t, time.Time{}, actGame.Actions[i].TimeStamp)
		}
		actGame.Actions[i].TimeStamp = time.Time{}
	}
	assert.Equal(t, expGame, actGame)
}

func TestDB(t *testing.T) {
	dbfs := map[dbName]persistence.DBFactory{
		memoryDB: memory.NewFactory(),
	}

	if !testing.Short() {
		// We assume you have mongodb stood up locally when running without -short
		// we change the uri because github actions set up a different mongodb replica set than run-rs does
		mongo, err := mongodb.NewFactory(`mongodb://127.0.0.1:27017,127.0.0.1:27018/?replicaSet=testReplSet`)
		require.NoError(t, err)

		dbfs[mongoDB] = mongo

		// We further assume you have mysql stood up locally when running without -short
		cfg := mysql.GetTestConfig()
		mySQLDB, err := mysql.NewFactory(context.Background(), cfg)
		if err != nil {
			t.Logf("Expected to connect, but got error: %q. This is expected when running locally.", err.Error())
			// if we got an error trying to connect, let's fallback to trying to connect to localhost's mysql
			cfg = mysql.GetTestConfigForLocal()
			mySQLDB, err = mysql.NewFactory(context.Background(), cfg)
		}
		require.NoError(t, err)

		dbfs[mysqlDB] = mySQLDB
	}

	for dbName, dbf := range dbfs {
		for testName, testFn := range tests {
			db, err := dbf.New(context.Background())
			require.NoError(t, err, string(dbName)+`:`+testName)
			t.Run(string(dbName)+`:`+testName, func(t1 *testing.T) { testFn(t1, dbName, db) })
		}
	}
}

func testCreatePlayersWithSimilarNames(t *testing.T, name dbName, db persistence.DB) {
	p1 := model.Player{
		ID:    model.PlayerID(`alice`),
		Name:  `alice`,
		Games: map[model.GameID]model.PlayerColor{},
	}

	assert.NoError(t, db.CreatePlayer(p1))

	p2 := model.Player{
		ID:    model.PlayerID(`Alice`),
		Name:  `Alice`,
		Games: map[model.GameID]model.PlayerColor{},
	}

	assert.NotEqual(t, p1.ID, p2.ID)
	assert.NoError(t, db.CreatePlayer(p2))
}

func testCreatePlayer(t *testing.T, name dbName, db persistence.DB) {
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
	expP2 := p2
	// Don't keep the same memory space for the games copy
	expP2.Games = map[model.GameID]model.PlayerColor{}

	assert.NoError(t, db.CreatePlayer(p2))

	actP2, err := db.GetPlayer(p2.ID)
	require.NoError(t, err)
	assert.Equal(t, expP2, actP2)

	alice, _, _ := testutils.EmptyAliceAndBob()
	assert.NoError(t, db.CreatePlayer(alice))

	// this is just a stub to allow us to add colors for the game
	g1 := model.Game{
		ID:              model.GameID(rand.Intn(1000)),
		Players:         []model.Player{alice, p2},
		BlockingPlayers: map[model.PlayerID]model.Blocker{p2.ID: model.PegCard},
		CurrentDealer:   alice.ID,
		PlayerColors:    nil,
		CurrentScores:   map[model.PlayerColor]int{},
		LagScores:       map[model.PlayerColor]int{},
		Phase:           model.Pegging,
		Hands:           map[model.PlayerID][]model.Card{},
		CutCard:         model.NewCardFromString(`KH`),
		Crib:            []model.Card{},
		PeggedCards:     make([]model.PeggedCard, 0, 8),
		Actions:         []model.PlayerAction{},
	}
	require.NoError(t, db.CreateGame(g1))

	require.NoError(t, db.AddPlayerColorToGame(p2.ID, model.Blue, g1.ID))

	expP2.Games = map[model.GameID]model.PlayerColor{
		g1.ID: model.Blue,
	}
	actP2, err = db.GetPlayer(p2.ID)
	require.NoError(t, err)
	assert.Equal(t, expP2, actP2)

	g2 := g1
	g2.PlayerColors = nil
	g2.ID = model.GameID(rand.Intn(1000))
	require.NoError(t, db.CreateGame(g2))

	require.NoError(t, db.AddPlayerColorToGame(p2.ID, model.Red, g2.ID))

	expP2.Games[g2.ID] = model.Red
	actP2, err = db.GetPlayer(p2.ID)
	require.NoError(t, err)
	assert.Equal(t, expP2, actP2)
}

func testCreateGame(t *testing.T, name dbName, db persistence.DB) {
	alice, bob, _ := testutils.EmptyAliceAndBob()

	g1 := model.Game{
		ID:              model.NewGameID(),
		Players:         []model.Player{alice, bob},
		BlockingPlayers: map[model.PlayerID]model.Blocker{bob.ID: model.PegCard},
		CurrentDealer:   alice.ID,
		PlayerColors:    map[model.PlayerID]model.PlayerColor{},
		CurrentScores:   map[model.PlayerColor]int{},
		LagScores:       map[model.PlayerColor]int{},
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
		Actions:     []model.PlayerAction{},
	}
	g1Copy := g1

	for i, p := range g1.Players {
		require.NoError(t, db.CreatePlayer(p))
		if c, ok := g1.PlayerColors[p.ID]; ok {
			g1.Players[i].Games = map[model.GameID]model.PlayerColor{
				g1.ID: c,
			}
		} else if name != mongoDB {
			// mongo doesn't give us this nicety, so we're gonna ignore it
			g1.Players[i].Games = map[model.GameID]model.PlayerColor{
				g1.ID: model.UnsetColor,
			}
		}
	}

	require.NoError(t, db.CreateGame(g1))

	actGame, err := db.GetGame(g1.ID)
	require.NoError(t, err, `expected to find game with id "%d"`, g1.ID)
	assert.Equal(t, g1Copy, actGame)
}

func testSaveGameMultipleTimes(t *testing.T, name dbName, db persistence.DB) {
	alice, bob, abAPIs := testutils.EmptyAliceAndBob()

	g, err := play.CreateGame([]model.Player{alice, bob}, abAPIs)
	require.NoError(t, err)

	for i, p := range g.Players {
		require.NoError(t, db.CreatePlayer(p))
		if c, ok := g.PlayerColors[p.ID]; ok {
			g.Players[i].Games = map[model.GameID]model.PlayerColor{
				g.ID: c,
			}
		}
	}

	_, err = db.GetGame(g.ID)
	require.Error(t, err)
	assert.EqualError(t, err, persistence.ErrGameNotFound.Error())

	var gCopy model.Game
	persistenceGameCopy(&gCopy, g)

	require.NoError(t, db.CreateGame(g))

	checkPersistedGame(t, name, db, gCopy)

	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        alice.ID,
		GameID:    g.ID,
		Overcomes: model.DealCards,
		Action:    model.DealAction{NumShuffles: 10},
	}, abAPIs))
	persistenceGameCopy(&gCopy, g)

	require.NoError(t, db.SaveGame(g))
	checkPersistedGame(t, name, db, gCopy)

	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        alice.ID,
		GameID:    g.ID,
		Overcomes: model.CribCard,
		Action:    model.BuildCribAction{Cards: []model.Card{g.Hands[alice.ID][0], g.Hands[alice.ID][1]}},
	}, abAPIs))
	persistenceGameCopy(&gCopy, g)

	require.NoError(t, db.SaveGame(g))
	checkPersistedGame(t, name, db, gCopy)

	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        bob.ID,
		GameID:    g.ID,
		Overcomes: model.CribCard,
		Action:    model.BuildCribAction{Cards: []model.Card{g.Hands[bob.ID][0], g.Hands[bob.ID][1]}},
	}, abAPIs))
	persistenceGameCopy(&gCopy, g)

	require.NoError(t, db.SaveGame(g))
	checkPersistedGame(t, name, db, gCopy)
}

func testSaveInteraction(t *testing.T, name dbName, db persistence.DB) {
	p1 := interaction.PlayerMeans{
		PlayerID:      model.PlayerID(rand.String(50)),
		PreferredMode: interaction.Localhost,
		Interactions: []interaction.Means{{
			Mode: interaction.Localhost,
			Info: string(`8383`),
		}},
	}
	p1Copy := p1

	require.NoError(t, db.CreatePlayer(model.Player{
		ID:   p1.PlayerID,
		Name: `testSaveInteractionStubPlayer`,
	}))

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
			Info: `8484`,
		}},
	}
	assert.NoError(t, db.SaveInteraction(p1update))

	actPM, err = db.GetInteraction(p1.PlayerID)
	assert.NoError(t, err)
	assert.NotEqual(t, p1Copy, actPM)
}

func testAddPlayerColorToGame(t *testing.T, name dbName, db persistence.DB) {
	alice, bob, abAPIs := testutils.EmptyAliceAndBob()

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

func testSaveGameWithMissingAction(t *testing.T, name dbName, db persistence.DB) {
	alice, bob, abAPIs := testutils.EmptyAliceAndBob()

	g, err := play.CreateGame([]model.Player{alice, bob}, abAPIs)
	require.NoError(t, err)
	for i, p := range g.Players {
		require.NoError(t, db.CreatePlayer(p))
		if c, ok := g.PlayerColors[p.ID]; ok {
			g.Players[i].Games = map[model.GameID]model.PlayerColor{
				g.ID: c,
			}
		}
	}

	_, err = db.GetGame(g.ID)
	require.Error(t, err)
	assert.EqualError(t, err, persistence.ErrGameNotFound.Error())

	var gCopy model.Game
	persistenceGameCopy(&gCopy, g)

	require.NoError(t, db.CreateGame(g))

	checkPersistedGame(t, name, db, gCopy)

	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        alice.ID,
		GameID:    g.ID,
		Overcomes: model.DealCards,
		Action:    model.DealAction{NumShuffles: 10},
	}, abAPIs))
	persistenceGameCopy(&gCopy, g)

	require.NoError(t, db.SaveGame(g))
	checkPersistedGame(t, name, db, gCopy)

	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        alice.ID,
		GameID:    g.ID,
		Overcomes: model.CribCard,
		Action:    model.BuildCribAction{Cards: []model.Card{g.Hands[alice.ID][0], g.Hands[alice.ID][1]}},
	}, abAPIs))
	persistenceGameCopy(&gCopy, g)

	require.NoError(t, db.SaveGame(g))
	checkPersistedGame(t, name, db, gCopy)

	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        bob.ID,
		GameID:    g.ID,
		Overcomes: model.CribCard,
		Action:    model.BuildCribAction{Cards: []model.Card{g.Hands[bob.ID][0], g.Hands[bob.ID][1]}},
	}, abAPIs))
	persistenceGameCopy(&gCopy, g)

	require.NoError(t, db.SaveGame(g))
	checkPersistedGame(t, name, db, gCopy)

	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        bob.ID,
		GameID:    g.ID,
		Overcomes: model.CutCard,
		Action:    model.CutDeckAction{Percentage: 0.600},
	}, abAPIs))
	persistenceGameCopy(&gCopy, g)

	// corrupt the "current" action
	i := len(g.Actions) - 1
	prevAction := g.Actions[i]
	badAction := prevAction
	badAction.ID = model.PlayerID(`nefario`)
	g.Actions[i] = badAction
	require.Error(t, db.SaveGame(g), `saving a game with an action by a player outside of it is a :badtime:`)

	badAction = prevAction
	badAction.GameID = g.ID + 1
	g.Actions[i] = badAction
	require.Error(t, db.SaveGame(g), `saving a game with an action on a different game is a :badtime:`)
	// set the latest action back to what it's supposed to be
	g.Actions[i] = prevAction

	// splice out an action
	prevActionSlice := g.Actions
	g.Actions = append(g.Actions[:1], g.Actions[2:]...)
	require.Error(t, db.SaveGame(g), `saving a game with a missing action is a :badtime:`)
	// set the action slice back to what it was
	g.Actions = prevActionSlice

	// corrupt a previous action
	badAction = g.Actions[1]
	badAction.Action = model.CountCribAction{Pts: 100}
	badAction.Overcomes = model.CountCrib
	g.Actions[1] = badAction
	if name == mysqlDB {
		// mysql is just storing one action per save. the previous ones can be corrupt as all get out
		// but as long as the latest one is fine, so are we
		assert.NoError(t, db.SaveGame(g), `saving a game with a corrupted action is a :badtime:`)
	} else {
		// this is because the noSQL databases are persisting ALL of the actions _every_ time
		assert.Error(t, db.SaveGame(g), `saving a game with a corrupted action is a :badtime:`)
	}
}

func TestTransactionality(t *testing.T) {
	// Right now (and probably ever) the memory persistence isn't transactional
	dbfs := map[dbName]persistence.DBFactory{
		// `memory`: func() persistence.DB { return inmem },
	}

	if !testing.Short() {
		// We assume you have mongodb stood up locally when running without -short
		// we change the uri because github actions set up a different mongodb replica set than run-rs does
		mongo, err := mongodb.NewFactory(`mongodb://127.0.0.1:27017,127.0.0.1:27018/?replicaSet=testReplSet`)
		require.NoError(t, err)

		dbfs[mongoDB] = mongo

		// We further assume you have mysql stood up locally when running without -short
		cfg := mysql.GetTestConfig()
		mySQLDB, err := mysql.NewFactory(context.Background(), cfg)
		if err != nil {
			t.Logf("Expected to connect, but got error: %q. This is expected when running locally.", err.Error())
			// if we got an error trying to connect, let's fallback to trying to connect to localhost's mysql
			cfg = mysql.GetTestConfigForLocal()
			mySQLDB, err = mysql.NewFactory(context.Background(), cfg)
		}
		require.NoError(t, err)

		dbfs[mysqlDB] = mySQLDB
	}

	txTests := map[string]txTest{
		`player service`:              playerTxTest,
		`game service`:                gameTxTest,
		`rollback the player service`: rollbackPlayerTxTest,
	}

	for databaseName, dbf := range dbfs {
		for testName, txTest := range txTests {
			db1, err := dbf.New(context.Background())
			require.NoError(t, err, string(databaseName)+`:`+testName)

			db2, err := dbf.New(context.Background())
			require.NoError(t, err, string(databaseName)+`:`+testName)

			db3, err := dbf.New(context.Background())
			require.NoError(t, err, string(databaseName)+`:`+testName)

			t.Run(string(databaseName)+`:`+testName, func(t1 *testing.T) { txTest(t1, databaseName, db1, db2, db3) })
		}
	}
}

type txTest func(t *testing.T, databaseName dbName, db1, db2, postCommitDB persistence.DB)

func playerTxTest(t *testing.T, databaseName dbName, db1, db2, postCommitDB persistence.DB) {
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
	err := db2.CreatePlayer(p1Mod)
	if databaseName == mysqlDB {
		assert.Error(t, err)
		assert.True(t, mysql.IsLockWaitTimeout(err))
	} else {
		assert.NoError(t, err)
	}

	err = db1.CreatePlayer(p1)
	assert.EqualError(t, err, persistence.ErrPlayerAlreadyExists.Error())

	savedP1, err := db1.GetPlayer(p1.ID)
	require.NoError(t, err)
	assert.Equal(t, p1, savedP1)

	savedP1Mod, err := db2.GetPlayer(p1.ID)
	if databaseName == mysqlDB {
		assert.Error(t, err)
		assert.EqualError(t, err, persistence.ErrPlayerNotFound.Error())
	} else {
		require.NoError(t, err)
		assert.NotEqual(t, p1, savedP1Mod)
	}

	assert.NoError(t, db1.Commit())

	postCommitP1, err := postCommitDB.GetPlayer(p1.ID)
	require.NoError(t, err)
	assert.Equal(t, p1, postCommitP1)
	assert.NotEqual(t, p1Mod, postCommitP1)
}

func rollbackPlayerTxTest(t *testing.T, databaseName dbName, db1, db2, postCommitDB persistence.DB) {
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

func gameTxTest(t *testing.T, databaseName dbName, db1, db2, postCommitDB persistence.DB) {
	alice, bob, _ := testutils.EmptyAliceAndBob()

	assert.NoError(t, db1.CreatePlayer(alice))
	assert.NoError(t, db1.CreatePlayer(bob))

	require.NoError(t, db1.Start())
	require.NoError(t, db2.Start())

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
	for i, p := range g1Copy.Players {
		if c, ok := g1.PlayerColors[p.ID]; ok {
			g1Copy.Players[i].Games = map[model.GameID]model.PlayerColor{
				g1.ID: c,
			}
		}
	}

	persistenceGameCopy(&g1Copy, g1)

	require.NoError(t, db1.CreateGame(g1))

	checkPersistedGame(t, databaseName, db1, g1Copy)

	actGame, err := db2.GetGame(g1.ID)
	assert.Error(t, err)
	assert.NotEqual(t, g1Copy, actGame)

	assert.NoError(t, db1.Commit())

	checkPersistedGame(t, databaseName, postCommitDB, g1Copy)
}
