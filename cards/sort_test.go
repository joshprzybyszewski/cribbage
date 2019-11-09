package cards

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortByValue(t *testing.T) {
	testCases := []struct {
		msg             string
		input           []string
		inputDescending bool
		expOutput       []string
	}{{
		msg:       `normal stuff`,
		input:     []string{`5s`, `2h`, `1s`},
		expOutput: []string{`1s`, `2h`, `5s`},
	}, {
		msg:       `normal stuff pt 2`,
		input:     []string{`ks`, `qh`, `js`, `10h`},
		expOutput: []string{`10h`, `js`, `qh`, `ks`},
	}, {
		msg:       `normal stuff pt 2 with duplicate`,
		input:     []string{`ks`, `qh`, `js`, `jh`, `10h`},
		expOutput: []string{`10h`, `js`, `jh`, `qh`, `ks`},
	}, {
		msg:             `normal stuff descending`,
		input:           []string{`1s`, `2h`, `5s`},
		inputDescending: true,
		expOutput:       []string{`5s`, `2h`, `1s`},
	}, {
		msg:             `normal stuff pt 2 descending with duplicate`,
		input:           []string{`ks`, `qh`, `js`, `jh`, `10h`},
		inputDescending: true,
		expOutput:       []string{`ks`, `qh`, `jh`, `js`, `10h`},
	}, {
		msg:             `normal stuff pt 2 descending`,
		input:           []string{`10h`, `js`, `qh`, `ks`},
		inputDescending: true,
		expOutput:       []string{`ks`, `qh`, `js`, `10h`},
	}}

	for _, tc := range testCases {
		inputCards := strToCards(tc.input)
		first := inputCards[0]
		actHand := SortByValue(inputCards, tc.inputDescending)
		for _, c := range actHand {
			assert.True(t, containsCard(tc.input, c), tc.msg+`: unexpected card `+c.String())
		}
		assert.Equal(t, strToCards(tc.expOutput), actHand)
		assert.Equal(t, first, inputCards[0], `check that the original list isn't touched`)
	}
}

func strToCards(s []string) []Card {
	c := make([]Card, len(s))
	for i, str := range s {
		c[i] = NewCardFromString(str)
	}
	return c
}

func containsCard(cs []string, c Card) bool {
	for _, cstr := range cs {
		if NewCardFromString(cstr).String() == c.String() {
			return true
		}
	}
	return false
}
