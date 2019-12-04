package persistence_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/persistence"
	"github.com/joshprzybyszewski/cribbage/server/persistence/memory"
	"github.com/joshprzybyszewski/cribbage/server/persistence/mongodb"
	"github.com/joshprzybyszewski/cribbage/utils/rand"
)

type dbTest func(*testing.T, persistence.DB)

var (
	tests = map[string]dbTest{
		`createPlayer`: testCreatePlayer,
		`saveGame`:     testSaveGame,
	}
)

func setup() (a, b model.Player, am, bm *interaction.Mock, pAPIs map[model.PlayerID]interaction.Player) {
	alice := model.Player{
		ID:   model.PlayerID(rand.String(100)),
		Name: `alice`,
	}
	bob := model.Player{
		ID:   model.PlayerID(rand.String(100)),
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

func TestDB(t *testing.T) {
	mongo, err := mongodb.New(``)
	require.NoError(t, err)

	dbs := map[string]persistence.DB{
		`memory`:  memory.New(),
		`mongodb`: mongo,
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
