package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/joshprzybyszewski/cribbage/model"
	"github.com/joshprzybyszewski/cribbage/utils/testutils"
)

func TestIsOver(t *testing.T) {
	testCases := []struct {
		msg     string
		game    model.Game
		expOver bool
	}{{
		msg: `easy detection`,
		game: model.Game{
			CurrentScores: map[model.PlayerColor]int{
				model.Blue: 121,
			},
		},
		expOver: true,
	}, {
		msg: `does not read lagging`,
		game: model.Game{
			LagScores: map[model.PlayerColor]int{
				model.Blue: 121,
			},
		},
		expOver: false,
	}, {
		msg: `does not win when close to winning`,
		game: model.Game{
			CurrentScores: map[model.PlayerColor]int{
				model.Blue: 120,
				model.Red:  120,
			},
		},
		expOver: false,
	}, {
		msg: `when green wins (what)`,
		game: model.Game{
			CurrentScores: map[model.PlayerColor]int{
				model.Blue:  120,
				model.Red:   120,
				model.Green: 121,
			},
		},
		expOver: true,
	}, {
		msg: `when we get more than max`,
		game: model.Game{
			CurrentScores: map[model.PlayerColor]int{
				model.Blue: 9001,
			},
		},
		expOver: true,
	}}

	for _, tc := range testCases {
		assert.Equal(t, tc.expOver, tc.game.IsOver(), tc.msg)
	}
}

func TestNumActions(t *testing.T) {
	alice, bob, charlie, diane := testutils.AliceBobCharlieDiane()

	testCases := []struct {
		msg    string
		game   model.Game
		expNum int
	}{{
		msg:    `no actions is fine`,
		game:   model.Game{},
		expNum: 0,
	}, {
		msg: `just returns a count of the actions`,
		game: model.Game{
			Actions: []model.PlayerAction{{
				ID:        alice.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{Card: model.NewCardFromString(`10s`)},
			}, {
				ID:        bob.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{Card: model.NewCardFromString(`10c`)},
			}, {
				ID:        charlie.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{Card: model.NewCardFromString(`10d`)},
			}, {
				ID:        diane.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{SayGo: true},
			}, {
				ID:        alice.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{SayGo: true},
			}, {
				ID:        bob.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{SayGo: true},
			}, {
				ID:        charlie.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{SayGo: true},
			}},
		},
		expNum: 7,
	}}

	for _, tc := range testCases {
		assert.Equal(t, tc.expNum, tc.game.NumActions(), tc.msg)
	}
}

func TestAddAction(t *testing.T) {
	alice, bob, _, _ := testutils.AliceBobCharlieDiane()

	g := model.Game{}
	assert.Zero(t, g.NumActions())

	g.AddAction(model.PlayerAction{
		ID:        alice.ID,
		Overcomes: model.PegCard,
		Action:    model.PegAction{Card: model.NewCardFromString(`10s`)},
	})
	assert.Equal(t, 1, g.NumActions())

	g.AddAction(model.PlayerAction{
		ID:        bob.ID,
		Overcomes: model.PegCard,
		Action:    model.PegAction{Card: model.NewCardFromString(`10c`)},
	})
	assert.Equal(t, 2, g.NumActions())
}

func TestCurrentPeg(t *testing.T) {
	alice, bob, charlie, diane := testutils.AliceBobCharlieDiane()

	testCases := []struct {
		msg    string
		game   model.Game
		expPeg int
	}{{
		msg:    `no pegged cards`,
		game:   model.Game{},
		expPeg: 0,
	}, {
		msg: `one pegged card`,
		game: model.Game{
			PeggedCards: []model.PeggedCard{
				model.NewPeggedCardFromString(alice.ID, `4c`, 0),
			},
			Actions: []model.PlayerAction{{
				ID:        alice.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{Card: model.NewCardFromString(`4c`)},
			}},
		},
		expPeg: 4,
	}, {
		msg: `two pegged cards`,
		game: model.Game{
			PeggedCards: []model.PeggedCard{
				model.NewPeggedCardFromString(alice.ID, `4c`, 0),
				model.NewPeggedCardFromString(bob.ID, `7c`, 1),
			},
			Actions: []model.PlayerAction{{
				ID:        alice.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{Card: model.NewCardFromString(`4c`)},
			}, {
				ID:        bob.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{Card: model.NewCardFromString(`7c`)},
			}},
		},
		expPeg: 11,
	}, {
		msg: `three pegged cards`,
		game: model.Game{
			PeggedCards: []model.PeggedCard{
				model.NewPeggedCardFromString(alice.ID, `4c`, 0),
				model.NewPeggedCardFromString(bob.ID, `7c`, 1),
				model.NewPeggedCardFromString(alice.ID, `10c`, 2),
			},
			Actions: []model.PlayerAction{{
				ID:        alice.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{Card: model.NewCardFromString(`4c`)},
			}, {
				ID:        bob.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{Card: model.NewCardFromString(`7c`)},
			}, {
				ID:        alice.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{Card: model.NewCardFromString(`10c`)},
			}},
		},
		expPeg: 21,
	}, {
		msg: `four pegged cards`,
		game: model.Game{
			PeggedCards: []model.PeggedCard{
				model.NewPeggedCardFromString(alice.ID, `4c`, 0),
				model.NewPeggedCardFromString(bob.ID, `7c`, 1),
				model.NewPeggedCardFromString(alice.ID, `10c`, 2),
				model.NewPeggedCardFromString(bob.ID, `9c`, 3),
			},
			Actions: []model.PlayerAction{{
				ID:        alice.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{Card: model.NewCardFromString(`4c`)},
			}, {
				ID:        bob.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{Card: model.NewCardFromString(`7c`)},
			}, {
				ID:        alice.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{Card: model.NewCardFromString(`10c`)},
			}, {
				ID:        bob.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{Card: model.NewCardFromString(`9c`)},
			}},
		},
		expPeg: 30,
	}, {
		msg: `after one go`,
		game: model.Game{
			PeggedCards: []model.PeggedCard{
				model.NewPeggedCardFromString(alice.ID, `4c`, 0),
				model.NewPeggedCardFromString(bob.ID, `7c`, 1),
				model.NewPeggedCardFromString(alice.ID, `10c`, 2),
				model.NewPeggedCardFromString(bob.ID, `9c`, 3),
			},
			Actions: []model.PlayerAction{{
				ID:        alice.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{Card: model.NewCardFromString(`4c`)},
			}, {
				ID:        bob.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{Card: model.NewCardFromString(`7c`)},
			}, {
				ID:        alice.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{Card: model.NewCardFromString(`10c`)},
			}, {
				ID:        bob.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{Card: model.NewCardFromString(`9c`)},
			}, {
				ID:        alice.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{SayGo: true},
			}},
		},
		expPeg: 30,
	}, {
		msg: `after two go's should reset`,
		game: model.Game{
			PeggedCards: []model.PeggedCard{
				model.NewPeggedCardFromString(alice.ID, `4c`, 0),
				model.NewPeggedCardFromString(bob.ID, `7c`, 1),
				model.NewPeggedCardFromString(alice.ID, `10c`, 2),
				model.NewPeggedCardFromString(bob.ID, `9c`, 3),
			},
			Actions: []model.PlayerAction{{
				ID:        alice.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{Card: model.NewCardFromString(`4c`)},
			}, {
				ID:        bob.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{Card: model.NewCardFromString(`7c`)},
			}, {
				ID:        alice.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{Card: model.NewCardFromString(`10c`)},
			}, {
				ID:        bob.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{Card: model.NewCardFromString(`9c`)},
			}, {
				ID:        alice.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{SayGo: true},
			}, {
				ID:        bob.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{SayGo: true},
			}},
		},
		expPeg: 0,
	}, {
		msg: `with three players, and three go's should reset`,
		game: model.Game{
			PeggedCards: []model.PeggedCard{
				model.NewPeggedCardFromString(alice.ID, `10s`, 0),
				model.NewPeggedCardFromString(bob.ID, `10c`, 1),
				model.NewPeggedCardFromString(charlie.ID, `10d`, 2),
			},
			Actions: []model.PlayerAction{{
				ID:        alice.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{Card: model.NewCardFromString(`10s`)},
			}, {
				ID:        bob.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{Card: model.NewCardFromString(`10c`)},
			}, {
				ID:        charlie.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{Card: model.NewCardFromString(`10d`)},
			}, {
				ID:        alice.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{SayGo: true},
			}, {
				ID:        bob.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{SayGo: true},
			}, {
				ID:        charlie.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{SayGo: true},
			}},
		},
		expPeg: 0,
	}, {
		msg: `with four players, and three go's should NOT reset`,
		game: model.Game{
			PeggedCards: []model.PeggedCard{
				model.NewPeggedCardFromString(alice.ID, `10s`, 0),
				model.NewPeggedCardFromString(bob.ID, `10c`, 1),
				model.NewPeggedCardFromString(charlie.ID, `10d`, 2),
			},
			Actions: []model.PlayerAction{{
				ID:        alice.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{Card: model.NewCardFromString(`10s`)},
			}, {
				ID:        bob.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{Card: model.NewCardFromString(`10c`)},
			}, {
				ID:        charlie.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{Card: model.NewCardFromString(`10d`)},
			}, {
				ID:        diane.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{SayGo: true},
			}, {
				ID:        alice.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{SayGo: true},
			}, {
				ID:        bob.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{SayGo: true},
			}},
		},
		expPeg: 30,
	}, {
		msg: `with four players, and four go's should reset`,
		game: model.Game{
			PeggedCards: []model.PeggedCard{
				model.NewPeggedCardFromString(alice.ID, `10s`, 0),
				model.NewPeggedCardFromString(bob.ID, `10c`, 1),
				model.NewPeggedCardFromString(charlie.ID, `10d`, 2),
			},
			Actions: []model.PlayerAction{{
				ID:        alice.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{Card: model.NewCardFromString(`10s`)},
			}, {
				ID:        bob.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{Card: model.NewCardFromString(`10c`)},
			}, {
				ID:        charlie.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{Card: model.NewCardFromString(`10d`)},
			}, {
				ID:        diane.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{SayGo: true},
			}, {
				ID:        alice.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{SayGo: true},
			}, {
				ID:        bob.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{SayGo: true},
			}, {
				ID:        charlie.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{SayGo: true},
			}},
		},
		expPeg: 0,
	}, {
		msg: `with four players, and the last player who played says go, should reset`,
		game: model.Game{
			PeggedCards: []model.PeggedCard{
				model.NewPeggedCardFromString(alice.ID, `10s`, 0),
				model.NewPeggedCardFromString(bob.ID, `10c`, 1),
				model.NewPeggedCardFromString(charlie.ID, `10d`, 2),
			},
			Actions: []model.PlayerAction{{
				ID:        alice.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{Card: model.NewCardFromString(`10s`)},
			}, {
				ID:        bob.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{Card: model.NewCardFromString(`10c`)},
			}, {
				ID:        charlie.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{Card: model.NewCardFromString(`10d`)},
			}, {
				ID:        charlie.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{SayGo: true},
			}},
		},
		expPeg: 0,
	}, {
		msg: `one pegged card, but after the pegging phase`,
		game: model.Game{
			Phase: model.Counting,
			PeggedCards: []model.PeggedCard{
				model.NewPeggedCardFromString(alice.ID, `4c`, 0),
			},
			Actions: []model.PlayerAction{{
				ID:        alice.ID,
				Overcomes: model.PegCard,
				Action:    model.PegAction{Card: model.NewCardFromString(`4c`)},
			}},
		},
		expPeg: 0,
	}}

	for _, tc := range testCases {
		assert.Equal(t, tc.expPeg, tc.game.CurrentPeg(), tc.msg)
	}
}

func TestPhaseString(t *testing.T) {
	testCases := []struct {
		input model.Phase
		exp   string
	}{{
		input: model.Deal,
		exp:   `Deal`,
	}, {
		input: model.BuildCribReady,
		exp:   `BuildCribReady`,
	}, {
		input: model.BuildCrib,
		exp:   `BuildCrib`,
	}, {
		input: model.CutReady,
		exp:   `CutReady`,
	}, {
		input: model.Cut,
		exp:   `Cut`,
	}, {
		input: model.PeggingReady,
		exp:   `PeggingReady`,
	}, {
		input: model.Pegging,
		exp:   `Pegging`,
	}, {
		input: model.CountingReady,
		exp:   `CountingReady`,
	}, {
		input: model.Counting,
		exp:   `Counting`,
	}, {
		input: model.CribCountingReady,
		exp:   `CribCountingReady`,
	}, {
		input: model.CribCounting,
		exp:   `CribCounting`,
	}, {
		input: model.DealingReady,
		exp:   `DealingReady`,
	}}

	for _, tc := range testCases {
		assert.Equal(t, tc.exp, tc.input.String())
	}
}
