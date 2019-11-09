package strategy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGiveCribHighestPotential(t *testing.T) {
	testCases := []struct {
		msg          string
		inputDesired int
		inputHand    []string
		expHand      []string
	}{{
		msg:          `obvious case`,
		inputDesired: 2,
		inputHand:    []string{`5s`, `5c`, `6d`, `9h`, `2h`, `1s`},
		expHand:      []string{`5s`, `5c`},
	}, {
		msg:          `obvious case still passses when requesting one card`,
		inputDesired: 1,
		inputHand:    []string{`5s`, `2c`, `6d`, `9h`, `1s`},
		expHand:      []string{`5s`},
	}}

	for _, tc := range testCases {
		actHand := GiveCribHighestPotential(tc.inputDesired, strToCards(tc.inputHand))
		for _, c := range actHand {
			assert.True(t, containsCard(tc.inputHand, c), tc.msg+`: unexpected card `+c.String())
		}
		assert.Equal(t, strToCards(tc.expHand), actHand)
	}
}

func TestGiveCribLowestPotential(t *testing.T) {
	testCases := []struct {
		msg          string
		inputDesired int
		inputHand    []string
		expHand      []string
	}{{
		msg:          `obvious case`,
		inputDesired: 2,
		inputHand:    []string{`5s`, `5c`, `5d`, `5h`, `2h`, `1s`},
		expHand:      []string{`2h`, `1s`},
	}, {
		msg:          `obvious case still passses when requesting one card`,
		inputDesired: 1,
		inputHand:    []string{`5s`, `5c`, `5d`, `5h`, `1s`},
		expHand:      []string{`1s`},
	}}

	for _, tc := range testCases {
		actHand := GiveCribLowestPotential(tc.inputDesired, strToCards(tc.inputHand))
		for _, c := range actHand {
			assert.True(t, containsCard(tc.inputHand, c), tc.msg+`: unexpected card `+c.String())
		}
		assert.Equal(t, strToCards(tc.expHand), actHand)
	}
}
