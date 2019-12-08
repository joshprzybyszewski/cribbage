package mapbson

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/joshprzybyszewski/cribbage/jsonutils"
	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/interaction"
	"github.com/joshprzybyszewski/cribbage/server/play"
	"github.com/joshprzybyszewski/cribbage/utils/testutils"
)

func TestCustomRegistryModelPlayer(t *testing.T) {
	registry := CustomRegistry()

	testCases := []struct {
		msg       string
		input     model.Player
		expOutput model.Player
	}{{
		msg: `custom map with three entries`,
		input: model.Player{
			ID:   model.PlayerID(`joshy`),
			Name: `joshy squashy`,
			Games: map[model.GameID]model.PlayerColor{
				model.GameID(42): model.Blue,
				model.GameID(77): model.Red,
			},
		},
		expOutput: model.Player{
			ID:   model.PlayerID(`joshy`),
			Name: `joshy squashy`,
			Games: map[model.GameID]model.PlayerColor{
				model.GameID(42): model.Blue,
				model.GameID(77): model.Red,
			},
		},
	}}

	for _, tc := range testCases {
		data, err := bson.MarshalWithRegistry(registry, tc.input)
		require.NoError(t, err, tc.msg)

		actOutput := model.Player{}
		err = bson.UnmarshalWithRegistry(registry, data, &actOutput)
		require.NoError(t, err, tc.msg)
		assert.Equal(t, tc.expOutput, actOutput, tc.msg)
	}
}

func TestCustomRegistryModelGame(t *testing.T) {
	registry := CustomRegistry()
	alice, bob, pAPIs := testutils.EmptyAliceAndBob()

	gPeg := gameAtPegging(t, alice, bob, pAPIs)
	gPegOutput := gPeg
	gPegOutput.Deck = nil

	testCases := []struct {
		msg       string
		input     model.Game
		expOutput model.Game
	}{{
		msg: `just a "normal" game`,
		input: model.Game{
			ID:      model.GameID(42),
			Players: []model.Player{alice, bob},
			PlayerColors: map[model.PlayerID]model.PlayerColor{
				alice.ID: model.Blue,
				bob.ID:   model.Red,
			},
			CurrentScores: map[model.PlayerColor]int{model.Blue: 55,
				model.Red: 64,
			},
			LagScores: map[model.PlayerColor]int{
				model.Red:  60,
				model.Blue: 34,
			},
			Phase: model.Pegging,
			BlockingPlayers: map[model.PlayerID]model.Blocker{
				alice.ID: model.PegCard,
			},
			CurrentDealer: bob.ID,
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
			Crib: []model.Card{
				model.NewCardFromString(`jh`),
				model.NewCardFromString(`jd`),
				model.NewCardFromString(`jc`),
				model.NewCardFromString(`js`),
			},
			CutCard: model.NewCardFromString(`5s`),
			PeggedCards: []model.PeggedCard{
				model.NewPeggedCardFromString(bob.ID, `4c`, 1),
			},
			Actions: []model.PlayerAction{{
				GameID:    model.GameID(42),
				ID:        bob.ID,
				Overcomes: model.PegCard,
				Action: model.PegAction{
					Card: model.NewCardFromString(`4c`),
				},
			}},
			Deck: model.NewDeck(),
		},
		expOutput: model.Game{
			ID:      model.GameID(42),
			Players: []model.Player{alice, bob},
			PlayerColors: map[model.PlayerID]model.PlayerColor{
				alice.ID: model.Blue,
				bob.ID:   model.Red,
			},
			CurrentScores: map[model.PlayerColor]int{model.Blue: 55,
				model.Red: 64,
			},
			LagScores: map[model.PlayerColor]int{
				model.Red:  60,
				model.Blue: 34,
			},
			Phase: model.Pegging,
			BlockingPlayers: map[model.PlayerID]model.Blocker{
				alice.ID: model.PegCard,
			},
			CurrentDealer: bob.ID,
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
			Crib: []model.Card{
				model.NewCardFromString(`jh`),
				model.NewCardFromString(`jd`),
				model.NewCardFromString(`jc`),
				model.NewCardFromString(`js`),
			},
			CutCard: model.NewCardFromString(`5s`),
			PeggedCards: []model.PeggedCard{
				model.NewPeggedCardFromString(bob.ID, `4c`, 1),
			},
			Actions: []model.PlayerAction{{
				GameID:    model.GameID(42),
				ID:        bob.ID,
				Overcomes: model.PegCard,
				Action: model.PegAction{
					Card: model.NewCardFromString(`4c`),
				},
			}},
			Deck: nil,
		},
	}, {
		msg:       `game at pegging`,
		input:     gPeg,
		expOutput: gPegOutput,
	}}

	for _, tc := range testCases {
		data1, err := bson.MarshalWithRegistry(registry, tc.input)
		require.NoError(t, err, tc.msg)
		data2 := make([]byte, len(data1))
		copy(data2, data1)

		actOutput := model.Game{}
		err = bson.UnmarshalWithRegistry(registry, data1, &actOutput)
		require.NoError(t, err, tc.msg)
		actOutput.Actions = nil
		expOutputNilled := tc.expOutput
		expOutputNilled.Actions = nil
		assert.Equal(t, expOutputNilled, actOutput, tc.msg)

		// Test deserialize from DB-BSON into model.Game
		tempGame := bson.M{}
		err = bson.UnmarshalWithRegistry(registry, data2, &tempGame)
		require.NoError(t, err, tc.msg)
		tempGameJSON, err := json.Marshal(tempGame)
		require.NoError(t, err, tc.msg)
		actOutput, err = jsonutils.UnmarshalGame(tempGameJSON)
		require.NoError(t, err, tc.msg)
		assert.Equal(t, tc.expOutput, actOutput, tc.msg)
	}
}

func gameAtPegging(t *testing.T, alice, bob model.Player, pAPIs map[model.PlayerID]interaction.Player) model.Game {
	g, err := play.CreateGame([]model.Player{alice, bob}, pAPIs)
	require.NoError(t, err)

	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        alice.ID,
		GameID:    g.ID,
		Overcomes: model.DealCards,
		Action:    model.DealAction{NumShuffles: 10},
	}, pAPIs))
	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        alice.ID,
		GameID:    g.ID,
		Overcomes: model.CribCard,
		Action:    model.BuildCribAction{Cards: []model.Card{g.Hands[alice.ID][0], g.Hands[alice.ID][1]}},
	}, pAPIs))
	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        bob.ID,
		GameID:    g.ID,
		Overcomes: model.CribCard,
		Action:    model.BuildCribAction{Cards: []model.Card{g.Hands[bob.ID][0], g.Hands[bob.ID][1]}},
	}, pAPIs))
	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        bob.ID,
		GameID:    g.ID,
		Overcomes: model.CutCard,
		Action:    model.CutDeckAction{Percentage: 0.314},
	}, pAPIs))
	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        bob.ID,
		GameID:    g.ID,
		Overcomes: model.PegCard,
		Action:    model.PegAction{Card: g.Hands[bob.ID][0]},
	}, pAPIs))

	return g
}
