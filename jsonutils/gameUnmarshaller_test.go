package jsonutils

import (
	"encoding/json"
	"testing"

	"github.com/joshprzybyszewski/cribbage/logic/scorer"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/server/play"
	"github.com/joshprzybyszewski/cribbage/utils/testutils"
)

func jsonCopyGame(input model.Game) model.Game {
	output := input

	if len(output.PlayerColors) == 0 {
		output.PlayerColors = nil
	}
	if len(output.BlockingPlayers) == 0 {
		output.BlockingPlayers = nil
	}
	if len(output.Hands) == 0 {
		output.Hands = nil
	}
	if len(output.Crib) == 0 {
		output.Crib = nil
	}
	if len(output.PeggedCards) == 0 {
		output.PeggedCards = nil
	}

	return output
}

func TestUnmarshalGame(t *testing.T) {
	alice, bob, _ := testutils.EmptyAliceAndBob()

	g5 := model.GameID(5)

	testCases := []struct {
		msg  string
		game model.Game
	}{{
		msg: `deal`,
		game: model.Game{
			ID:              g5,
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
			Actions: []model.PlayerAction{{
				GameID:    g5,
				ID:        alice.ID,
				Overcomes: model.DealCards,
				Action: model.DealAction{
					NumShuffles: 543,
				},
			}},
		},
	}}

	for _, tc := range testCases {
		checkMarshalUnmarshal(t, tc.game, tc.msg)
	}
}

func checkMarshalUnmarshal(t *testing.T, g model.Game, msg string) {
	gameCopy := jsonCopyGame(g)

	b, err := json.Marshal(g)
	require.NoError(t, err, msg)

	actGame, err := UnmarshalGame(b)
	require.NoError(t, err, msg)

	if gameCopy.PlayerColors == nil {
		assert.NotNil(t, actGame.PlayerColors)
		actGame.PlayerColors = nil
	}
	if gameCopy.BlockingPlayers == nil {
		assert.NotNil(t, actGame.BlockingPlayers)
		actGame.BlockingPlayers = nil
	}
	if gameCopy.Hands == nil {
		assert.NotNil(t, actGame.Hands)
		actGame.Hands = nil
	}

	assert.Equal(t, gameCopy, actGame, msg)
}

func TestGameAtAllStages(t *testing.T) {
	alice, bob, pAPIs := testutils.EmptyAliceAndBob()
	g, err := play.CreateGame([]model.Player{alice, bob}, pAPIs)
	require.NoError(t, err)
	checkMarshalUnmarshal(t, g, `after creation`)

	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        alice.ID,
		GameID:    g.ID,
		Overcomes: model.DealCards,
		Action:    model.DealAction{NumShuffles: 10},
	}, pAPIs))
	checkMarshalUnmarshal(t, g, `after deal`)

	g.Hands[alice.ID] = []model.Card{
		model.NewCardFromString(`5s`),
		model.NewCardFromString(`5c`),
		model.NewCardFromString(`1s`),
		model.NewCardFromString(`2s`),
		model.NewCardFromString(`3s`),
		model.NewCardFromString(`4s`),
	}
	g.Hands[bob.ID] = []model.Card{
		model.NewCardFromString(`5d`),
		model.NewCardFromString(`5h`),
		model.NewCardFromString(`1c`),
		model.NewCardFromString(`2c`),
		model.NewCardFromString(`3c`),
		model.NewCardFromString(`4c`),
	}
	checkMarshalUnmarshal(t, g, `after deal`)

	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        alice.ID,
		GameID:    g.ID,
		Overcomes: model.CribCard,
		Action:    model.BuildCribAction{Cards: []model.Card{g.Hands[alice.ID][0], g.Hands[alice.ID][1]}},
	}, pAPIs))
	checkMarshalUnmarshal(t, g, `after crib from alice`)

	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        bob.ID,
		GameID:    g.ID,
		Overcomes: model.CribCard,
		Action:    model.BuildCribAction{Cards: []model.Card{g.Hands[bob.ID][0], g.Hands[bob.ID][1]}},
	}, pAPIs))
	checkMarshalUnmarshal(t, g, `after crib from bob`)

	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        bob.ID,
		GameID:    g.ID,
		Overcomes: model.CutCard,
		Action:    model.CutDeckAction{Percentage: 0.314},
	}, pAPIs))
	checkMarshalUnmarshal(t, g, `after cut from bob`)

	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        bob.ID,
		GameID:    g.ID,
		Overcomes: model.PegCard,
		Action:    model.PegAction{Card: g.Hands[bob.ID][0]},
	}, pAPIs))
	checkMarshalUnmarshal(t, g, `after bob pegs`)

	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        alice.ID,
		GameID:    g.ID,
		Overcomes: model.PegCard,
		Action:    model.PegAction{Card: g.Hands[alice.ID][0]},
	}, pAPIs))
	checkMarshalUnmarshal(t, g, `after alice pegs`)

	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        bob.ID,
		GameID:    g.ID,
		Overcomes: model.PegCard,
		Action:    model.PegAction{Card: g.Hands[bob.ID][1]},
	}, pAPIs))
	checkMarshalUnmarshal(t, g, `after bob pegs`)

	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        alice.ID,
		GameID:    g.ID,
		Overcomes: model.PegCard,
		Action:    model.PegAction{Card: g.Hands[alice.ID][1]},
	}, pAPIs))
	checkMarshalUnmarshal(t, g, `after alice pegs`)

	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        bob.ID,
		GameID:    g.ID,
		Overcomes: model.PegCard,
		Action:    model.PegAction{Card: g.Hands[bob.ID][2]},
	}, pAPIs))
	checkMarshalUnmarshal(t, g, `after bob pegs`)

	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        alice.ID,
		GameID:    g.ID,
		Overcomes: model.PegCard,
		Action:    model.PegAction{Card: g.Hands[alice.ID][2]},
	}, pAPIs))
	checkMarshalUnmarshal(t, g, `after alice pegs`)

	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        bob.ID,
		GameID:    g.ID,
		Overcomes: model.PegCard,
		Action:    model.PegAction{Card: g.Hands[bob.ID][3]},
	}, pAPIs))
	checkMarshalUnmarshal(t, g, `after bob pegs`)

	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        alice.ID,
		GameID:    g.ID,
		Overcomes: model.PegCard,
		Action:    model.PegAction{Card: g.Hands[alice.ID][3]},
	}, pAPIs))
	checkMarshalUnmarshal(t, g, `after alice pegs`)

	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        bob.ID,
		GameID:    g.ID,
		Overcomes: model.CountHand,
		Action:    model.CountHandAction{Pts: scorer.HandPoints(g.CutCard, g.Hands[bob.ID])},
	}, pAPIs))
	checkMarshalUnmarshal(t, g, `after bob scores`)

	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        alice.ID,
		GameID:    g.ID,
		Overcomes: model.CountHand,
		Action:    model.CountHandAction{Pts: scorer.HandPoints(g.CutCard, g.Hands[alice.ID])},
	}, pAPIs))
	checkMarshalUnmarshal(t, g, `after alice scores`)

	require.NoError(t, play.HandleAction(&g, model.PlayerAction{
		ID:        alice.ID,
		GameID:    g.ID,
		Overcomes: model.CountCrib,
		Action:    model.CountCribAction{Pts: scorer.CribPoints(g.CutCard, g.Crib)},
	}, pAPIs))
	checkMarshalUnmarshal(t, g, `after alice scores crib`)
}
