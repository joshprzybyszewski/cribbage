package strategy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAvoidCribFifteens(t *testing.T) {
	testCases := []struct {
		msg          string
		inputDesired int
		inputHand    []string
		expHand      []string
	}{{
		msg:          `obvious case`,
		inputDesired: 2,
		inputHand:    []string{`5s`, `5c`, `5d`, `5h`, `2h`, `1s`},
		expHand:      []string{`1s`, `2h`},
	}}

	for _, tc := range testCases {
		actHand := AvoidCribFifteens(tc.inputDesired, strToCards(tc.inputHand))
		assert.Equal(t, strToCards(tc.expHand), actHand)
	}
}

func TestGiveCribFifteens(t *testing.T) {
	testCases := []struct {
		msg          string
		inputDesired int
		inputHand    []string
		expHand      []string
	}{{
		msg:          `obvious case`,
		inputDesired: 2,
		inputHand:    []string{`5s`, `10c`, `8d`, `9h`, `2h`, `1s`},
		expHand:      []string{`5s`, `10c`},
	}}

	for _, tc := range testCases {
		actHand := AvoidCribFifteens(tc.inputDesired, strToCards(tc.inputHand))
		assert.Equal(t, strToCards(tc.expHand), actHand)
	}
}
