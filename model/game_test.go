package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getFourPlayers() (alice, bob, charlie, diane Player) {
	alice = Player{
		ID:   PlayerID(`a`),
		Name: `alice`,
	}
	bob = Player{
		ID:   PlayerID(`b`),
		Name: `bob`,
	}
	charlie = Player{
		ID:   PlayerID(`c`),
		Name: `charlie`,
	}
	diane = Player{
		ID:   PlayerID(`d`),
		Name: `diane`,
	}
	return alice, bob, charlie, diane
}

func TestIsOver(t *testing.T) {
	testCases := []struct {
		msg     string
		game    Game
		expOver bool
	}{{
		msg: `easy detection`,
		game: Game{
			CurrentScores: map[PlayerColor]int{
				Blue: 121,
			},
		},
		expOver: true,
	}, {
		msg: `does not read lagging`,
		game: Game{
			LagScores: map[PlayerColor]int{
				Blue: 121,
			},
		},
		expOver: false,
	}, {
		msg: `does not win when close to winning`,
		game: Game{
			CurrentScores: map[PlayerColor]int{
				Blue: 120,
				Red:  120,
			},
		},
		expOver: false,
	}, {
		msg: `when green wins (what)`,
		game: Game{
			CurrentScores: map[PlayerColor]int{
				Blue:  120,
				Red:   120,
				Green: 121,
			},
		},
		expOver: true,
	}, {
		msg: `when we get more than max`,
		game: Game{
			CurrentScores: map[PlayerColor]int{
				Blue: 9001,
			},
		},
		expOver: true,
	}}

	for _, tc := range testCases {
		assert.Equal(t, tc.expOver, tc.game.IsOver(), tc.msg)
	}
}

func TestNumActions(t *testing.T) {
	alice, bob, charlie, diane := getFourPlayers()

	testCases := []struct {
		msg    string
		game   Game
		expNum int
	}{{
		msg:    `no actions is fine`,
		game:   Game{},
		expNum: 0,
	}, {
		msg: `just returns a count of the actions`,
		game: Game{
			actions: []PlayerAction{{
				ID:        alice.ID,
				Overcomes: PegCard,
				Action:    PegAction{Card: NewCardFromString(`10s`)},
			}, {
				ID:        bob.ID,
				Overcomes: PegCard,
				Action:    PegAction{Card: NewCardFromString(`10c`)},
			}, {
				ID:        charlie.ID,
				Overcomes: PegCard,
				Action:    PegAction{Card: NewCardFromString(`10d`)},
			}, {
				ID:        diane.ID,
				Overcomes: PegCard,
				Action:    PegAction{SayGo: true},
			}, {
				ID:        alice.ID,
				Overcomes: PegCard,
				Action:    PegAction{SayGo: true},
			}, {
				ID:        bob.ID,
				Overcomes: PegCard,
				Action:    PegAction{SayGo: true},
			}, {
				ID:        charlie.ID,
				Overcomes: PegCard,
				Action:    PegAction{SayGo: true},
			}},
		},
		expNum: 7,
	}}

	for _, tc := range testCases {
		assert.Equal(t, tc.expNum, tc.game.NumActions(), tc.msg)
	}
}

func TestAddAction(t *testing.T) {
	alice, bob, _, _ := getFourPlayers()

	g := Game{}
	assert.Zero(t, g.NumActions())

	g.AddAction(PlayerAction{
		ID:        alice.ID,
		Overcomes: PegCard,
		Action:    PegAction{Card: NewCardFromString(`10s`)},
	})
	assert.Equal(t, 1, g.NumActions())

	g.AddAction(PlayerAction{
		ID:        bob.ID,
		Overcomes: PegCard,
		Action:    PegAction{Card: NewCardFromString(`10c`)},
	})
	assert.Equal(t, 2, g.NumActions())
}

func TestCurrentPeg(t *testing.T) {
	alice, bob, charlie, diane := getFourPlayers()

	testCases := []struct {
		msg    string
		game   Game
		expPeg int
	}{{
		msg:    `no pegged cards`,
		game:   Game{},
		expPeg: 0,
	}, {
		msg: `one pegged card`,
		game: Game{
			PeggedCards: []PeggedCard{
				NewPeggedCardFromString(alice.ID, `4c`, 0),
			},
			actions: []PlayerAction{{
				ID:        alice.ID,
				Overcomes: PegCard,
				Action:    PegAction{Card: NewCardFromString(`4c`)},
			}},
		},
		expPeg: 4,
	}, {
		msg: `two pegged cards`,
		game: Game{
			PeggedCards: []PeggedCard{
				NewPeggedCardFromString(alice.ID, `4c`, 0),
				NewPeggedCardFromString(bob.ID, `7c`, 1),
			},
			actions: []PlayerAction{{
				ID:        alice.ID,
				Overcomes: PegCard,
				Action:    PegAction{Card: NewCardFromString(`4c`)},
			}, {
				ID:        bob.ID,
				Overcomes: PegCard,
				Action:    PegAction{Card: NewCardFromString(`7c`)},
			}},
		},
		expPeg: 11,
	}, {
		msg: `three pegged cards`,
		game: Game{
			PeggedCards: []PeggedCard{
				NewPeggedCardFromString(alice.ID, `4c`, 0),
				NewPeggedCardFromString(bob.ID, `7c`, 1),
				NewPeggedCardFromString(alice.ID, `10c`, 2),
			},
			actions: []PlayerAction{{
				ID:        alice.ID,
				Overcomes: PegCard,
				Action:    PegAction{Card: NewCardFromString(`4c`)},
			}, {
				ID:        bob.ID,
				Overcomes: PegCard,
				Action:    PegAction{Card: NewCardFromString(`7c`)},
			}, {
				ID:        alice.ID,
				Overcomes: PegCard,
				Action:    PegAction{Card: NewCardFromString(`10c`)},
			}},
		},
		expPeg: 21,
	}, {
		msg: `four pegged cards`,
		game: Game{
			PeggedCards: []PeggedCard{
				NewPeggedCardFromString(alice.ID, `4c`, 0),
				NewPeggedCardFromString(bob.ID, `7c`, 1),
				NewPeggedCardFromString(alice.ID, `10c`, 2),
				NewPeggedCardFromString(bob.ID, `9c`, 3),
			},
			actions: []PlayerAction{{
				ID:        alice.ID,
				Overcomes: PegCard,
				Action:    PegAction{Card: NewCardFromString(`4c`)},
			}, {
				ID:        bob.ID,
				Overcomes: PegCard,
				Action:    PegAction{Card: NewCardFromString(`7c`)},
			}, {
				ID:        alice.ID,
				Overcomes: PegCard,
				Action:    PegAction{Card: NewCardFromString(`10c`)},
			}, {
				ID:        bob.ID,
				Overcomes: PegCard,
				Action:    PegAction{Card: NewCardFromString(`9c`)},
			}},
		},
		expPeg: 30,
	}, {
		msg: `after one go`,
		game: Game{
			PeggedCards: []PeggedCard{
				NewPeggedCardFromString(alice.ID, `4c`, 0),
				NewPeggedCardFromString(bob.ID, `7c`, 1),
				NewPeggedCardFromString(alice.ID, `10c`, 2),
				NewPeggedCardFromString(bob.ID, `9c`, 3),
			},
			actions: []PlayerAction{{
				ID:        alice.ID,
				Overcomes: PegCard,
				Action:    PegAction{Card: NewCardFromString(`4c`)},
			}, {
				ID:        bob.ID,
				Overcomes: PegCard,
				Action:    PegAction{Card: NewCardFromString(`7c`)},
			}, {
				ID:        alice.ID,
				Overcomes: PegCard,
				Action:    PegAction{Card: NewCardFromString(`10c`)},
			}, {
				ID:        bob.ID,
				Overcomes: PegCard,
				Action:    PegAction{Card: NewCardFromString(`9c`)},
			}, {
				ID:        alice.ID,
				Overcomes: PegCard,
				Action:    PegAction{SayGo: true},
			}},
		},
		expPeg: 30,
	}, {
		msg: `after two go's should reset`,
		game: Game{
			PeggedCards: []PeggedCard{
				NewPeggedCardFromString(alice.ID, `4c`, 0),
				NewPeggedCardFromString(bob.ID, `7c`, 1),
				NewPeggedCardFromString(alice.ID, `10c`, 2),
				NewPeggedCardFromString(bob.ID, `9c`, 3),
			},
			actions: []PlayerAction{{
				ID:        alice.ID,
				Overcomes: PegCard,
				Action:    PegAction{Card: NewCardFromString(`4c`)},
			}, {
				ID:        bob.ID,
				Overcomes: PegCard,
				Action:    PegAction{Card: NewCardFromString(`7c`)},
			}, {
				ID:        alice.ID,
				Overcomes: PegCard,
				Action:    PegAction{Card: NewCardFromString(`10c`)},
			}, {
				ID:        bob.ID,
				Overcomes: PegCard,
				Action:    PegAction{Card: NewCardFromString(`9c`)},
			}, {
				ID:        alice.ID,
				Overcomes: PegCard,
				Action:    PegAction{SayGo: true},
			}, {
				ID:        bob.ID,
				Overcomes: PegCard,
				Action:    PegAction{SayGo: true},
			}},
		},
		expPeg: 0,
	}, {
		msg: `with three players, and three go's should reset`,
		game: Game{
			PeggedCards: []PeggedCard{
				NewPeggedCardFromString(alice.ID, `10s`, 0),
				NewPeggedCardFromString(bob.ID, `10c`, 1),
				NewPeggedCardFromString(charlie.ID, `10d`, 2),
			},
			actions: []PlayerAction{{
				ID:        alice.ID,
				Overcomes: PegCard,
				Action:    PegAction{Card: NewCardFromString(`10s`)},
			}, {
				ID:        bob.ID,
				Overcomes: PegCard,
				Action:    PegAction{Card: NewCardFromString(`10c`)},
			}, {
				ID:        charlie.ID,
				Overcomes: PegCard,
				Action:    PegAction{Card: NewCardFromString(`10d`)},
			}, {
				ID:        alice.ID,
				Overcomes: PegCard,
				Action:    PegAction{SayGo: true},
			}, {
				ID:        bob.ID,
				Overcomes: PegCard,
				Action:    PegAction{SayGo: true},
			}, {
				ID:        charlie.ID,
				Overcomes: PegCard,
				Action:    PegAction{SayGo: true},
			}},
		},
		expPeg: 0,
	}, {
		msg: `with four players, and three go's should NOT reset`,
		game: Game{
			PeggedCards: []PeggedCard{
				NewPeggedCardFromString(alice.ID, `10s`, 0),
				NewPeggedCardFromString(bob.ID, `10c`, 1),
				NewPeggedCardFromString(charlie.ID, `10d`, 2),
			},
			actions: []PlayerAction{{
				ID:        alice.ID,
				Overcomes: PegCard,
				Action:    PegAction{Card: NewCardFromString(`10s`)},
			}, {
				ID:        bob.ID,
				Overcomes: PegCard,
				Action:    PegAction{Card: NewCardFromString(`10c`)},
			}, {
				ID:        charlie.ID,
				Overcomes: PegCard,
				Action:    PegAction{Card: NewCardFromString(`10d`)},
			}, {
				ID:        diane.ID,
				Overcomes: PegCard,
				Action:    PegAction{SayGo: true},
			}, {
				ID:        alice.ID,
				Overcomes: PegCard,
				Action:    PegAction{SayGo: true},
			}, {
				ID:        bob.ID,
				Overcomes: PegCard,
				Action:    PegAction{SayGo: true},
			}},
		},
		expPeg: 30,
	}, {
		msg: `with four players, and four go's should reset`,
		game: Game{
			PeggedCards: []PeggedCard{
				NewPeggedCardFromString(alice.ID, `10s`, 0),
				NewPeggedCardFromString(bob.ID, `10c`, 1),
				NewPeggedCardFromString(charlie.ID, `10d`, 2),
			},
			actions: []PlayerAction{{
				ID:        alice.ID,
				Overcomes: PegCard,
				Action:    PegAction{Card: NewCardFromString(`10s`)},
			}, {
				ID:        bob.ID,
				Overcomes: PegCard,
				Action:    PegAction{Card: NewCardFromString(`10c`)},
			}, {
				ID:        charlie.ID,
				Overcomes: PegCard,
				Action:    PegAction{Card: NewCardFromString(`10d`)},
			}, {
				ID:        diane.ID,
				Overcomes: PegCard,
				Action:    PegAction{SayGo: true},
			}, {
				ID:        alice.ID,
				Overcomes: PegCard,
				Action:    PegAction{SayGo: true},
			}, {
				ID:        bob.ID,
				Overcomes: PegCard,
				Action:    PegAction{SayGo: true},
			}, {
				ID:        charlie.ID,
				Overcomes: PegCard,
				Action:    PegAction{SayGo: true},
			}},
		},
		expPeg: 0,
	}, {
		msg: `with four players, and the last player who played says go, should reset`,
		game: Game{
			PeggedCards: []PeggedCard{
				NewPeggedCardFromString(alice.ID, `10s`, 0),
				NewPeggedCardFromString(bob.ID, `10c`, 1),
				NewPeggedCardFromString(charlie.ID, `10d`, 2),
			},
			actions: []PlayerAction{{
				ID:        alice.ID,
				Overcomes: PegCard,
				Action:    PegAction{Card: NewCardFromString(`10s`)},
			}, {
				ID:        bob.ID,
				Overcomes: PegCard,
				Action:    PegAction{Card: NewCardFromString(`10c`)},
			}, {
				ID:        charlie.ID,
				Overcomes: PegCard,
				Action:    PegAction{Card: NewCardFromString(`10d`)},
			}, {
				ID:        charlie.ID,
				Overcomes: PegCard,
				Action:    PegAction{SayGo: true},
			}},
		},
		expPeg: 0,
	}}

	for _, tc := range testCases {
		assert.Equal(t, tc.expPeg, tc.game.CurrentPeg(), tc.msg)
	}
}
