package pegging

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/joshprzybyszewski/cribbage/cards"
)

func TestPointsForCard(t *testing.T) {
	testCases := []struct {
		msg        string
		inputCards []string
		inputCard  string
		expVal     int
		expErr     error
	}{{
		msg:        `no points`,
		inputCards: []string{`10C`},
		inputCard:  `4C`,
		expVal:     0,
		expErr:     nil,
	}, {
		msg:        `hits fifteen`,
		inputCards: []string{`10C`},
		inputCard:  `5C`,
		expVal:     2,
		expErr:     nil,
	}, {
		msg:        `hits 31`,
		inputCards: []string{`10C`, `10D`, `10H`},
		inputCard:  `1C`,
		expVal:     2,
		expErr:     nil,
	}, {
		msg:        `hits pair`,
		inputCards: []string{`10C`},
		inputCard:  `10D`,
		expVal:     2,
		expErr:     nil,
	}, {
		msg:        `hits triplet`,
		inputCards: []string{`10C`, `10D`},
		inputCard:  `10H`,
		expVal:     6,
		expErr:     nil,
	}, {
		msg:        `hits quad`,
		inputCards: []string{`2C`, `2D`, `2S`},
		inputCard:  `2H`,
		expVal:     12,
		expErr:     nil,
	}, {
		msg:        `run of three`,
		inputCards: []string{`1C`, `3D`},
		inputCard:  `2H`,
		expVal:     3,
		expErr:     nil,
	}, {
		msg:        `run of four`,
		inputCards: []string{`1C`, `2D`, `4H`},
		inputCard:  `3H`,
		expVal:     4,
		expErr:     nil,
	}, {
		msg:        `run of five`,
		inputCards: []string{`6C`, `2D`, `4H`, `5H`},
		inputCard:  `3H`,
		expVal:     5,
		expErr:     nil,
	}, {
		msg:        `run of 6`,
		inputCards: []string{`1C`, `2D`, `4H`, `5H`, `6H`},
		inputCard:  `3H`,
		expVal:     6,
		expErr:     nil,
	}, {
		msg:        `run of 7`,
		inputCards: []string{`1C`, `2D`, `4H`, `5H`, `6H`, `7H`},
		inputCard:  `3H`,
		expVal:     7,
		expErr:     nil,
	}, {
		msg:        `run with 15`,
		inputCards: []string{`1C`, `2D`, `4H`, `5H`},
		inputCard:  `3H`,
		expVal:     7,
		expErr:     nil,
	}, {
		msg:        `run with 31`,
		inputCards: []string{`10C`, `10D`, `2H`, `2C`, `4H`},
		inputCard:  `3H`,
		expVal:     5,
		expErr:     nil,
	}, {
		msg:        `run of three (after run of three)`,
		inputCards: []string{`1C`, `3D`, `2D`},
		inputCard:  `1H`,
		expVal:     3,
		expErr:     nil,
	}, {
		msg:        `run of three (after run of three)`,
		inputCards: []string{`3C`, `1D`, `2D`},
		inputCard:  `3H`,
		expVal:     3,
		expErr:     nil,
	}, {
		msg:        `close to a run, but isn't`,
		inputCards: []string{`3C`, `1D`, `10C`, `2D`},
		inputCard:  `4H`,
		expVal:     0,
		expErr:     nil,
	}, {
		msg:        `not a run, hits 31`,
		inputCards: []string{`4♥︎`, `8♥︎`, `7♠︎`, `6♣︎`, `5♥︎`},
		inputCard:  `AC`,
		expVal:     2,
		expErr:     nil,
	}, {
		msg:        `looks like a run, but over a 31`,
		inputCards: []string{`4♠︎`, `J♦`, `K♥︎`, `A♦`, `9♣︎`, `7♣︎`, `8♠︎`},
		inputCard:  `9D`,
		expVal:     0,
		expErr:     nil,
	}}

	for _, tc := range testCases {
		c := make([]cards.Card, len(tc.inputCards))
		for i, ic := range tc.inputCards {
			c[i] = cards.NewCardFromString(ic)
		}
		next := cards.NewCardFromString(tc.inputCard)
		actVal, actErr := PointsForCard(c, next)
		assert.Equal(t, tc.expErr, actErr, `unexpected error for test "%s"`, tc.msg)
		assert.Equal(t, tc.expVal, actVal, `unexpected value for test "%s"`, tc.msg)
	}
}

func TestPointsErrorCase(t *testing.T) {
	testCases := []struct {
		msg        string
		inputCards []string
		inputCard  string
		expVal     int
		expErr     error
	}{{
		msg:        `same card twice`,
		inputCards: []string{`10C`, `10C`},
		inputCard:  `4C`,
		expVal:     0,
		expErr:     errSameCardTwice,
	}, {
		msg:        `same card twice -- one in the input card`,
		inputCards: []string{`10C`},
		inputCard:  `10C`,
		expVal:     0,
		expErr:     errSameCardTwice,
	}, {
		msg:        `definitely too many cards`,
		inputCards: []string{`1S`, `1C`, `1D`, `1H`, `2S`, `2C`, `2D`, `2H`, `3S`, `3C`, `3D`, `3H`, `4S`, `4C`, `4D`, `4H`},
		inputCard:  `5S`,
		expVal:     0,
		expErr:     errTooManyCards,
	}, {
		msg:        `run of 8`,
		inputCards: []string{`1C`, `2D`, `4H`, `5H`, `6H`, `7H`, `3H`},
		inputCard:  `8H`,
		expVal:     0,
		expErr:     errRunTooLong,
	}}

	for _, tc := range testCases {
		c := make([]cards.Card, len(tc.inputCards))
		for i, ic := range tc.inputCards {
			c[i] = cards.NewCardFromString(ic)
		}
		next := cards.NewCardFromString(tc.inputCard)
		actVal, actErr := PointsForCard(c, next)
		assert.Equal(t, tc.expErr, actErr, `unexpected error for test "%s"`, tc.msg)
		assert.Equal(t, tc.expVal, actVal, `unexpected value for test "%s"`, tc.msg)
	}
}
