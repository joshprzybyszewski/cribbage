package cards

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandToInt32(t *testing.T) {
	testCases := []struct {
		msg       string
		inputHand []Card
		expInt32  int32
	}{{
		msg:       `not 4 cards`,
		inputHand: nil,
		expInt32:  -1,
	}, {
		msg:       `first 4 cards in the deck`,
		inputHand: []Card{NewCardFromString(`AS`), NewCardFromString(`2S`), NewCardFromString(`3S`), NewCardFromString(`4S`)},
		expInt32:  0x00010203,
	}, {
		msg:       `first 4 cards in the Hearts`,
		inputHand: []Card{NewCardFromString(`AH`), NewCardFromString(`2H`), NewCardFromString(`3H`), NewCardFromString(`4H`)},
		expInt32:  0x2728292a,
	}, {
		msg:       `last 4 cards in the deck`,
		inputHand: []Card{NewCardFromString(`10H`), NewCardFromString(`JH`), NewCardFromString(`QH`), NewCardFromString(`KH`)},
		expInt32:  0x30313233,
	}, {
		msg:       `2S, 3s,4s,5s not in order`,
		inputHand: []Card{NewCardFromString(`5s`), NewCardFromString(`2s`), NewCardFromString(`4s`), NewCardFromString(`3s`)},
		expInt32:  0x01020304,
	}}

	for _, tc := range testCases {
		assert.Equal(t, tc.expInt32, HandToInt32(tc.inputHand), tc.msg)
	}
}

func TestInt32ToHand(t *testing.T) {
	testCases := []struct {
		msg        string
		inputInt32 int32
		expHand    []Card
	}{{
		msg:        `first 4 cards in the deck`,
		inputInt32: 0x00010203,
		expHand:    []Card{NewCardFromString(`AS`), NewCardFromString(`2S`), NewCardFromString(`3S`), NewCardFromString(`4S`)},
	}, {
		msg:        `first 4 cards in the Hearts`,
		inputInt32: 0x2728292a,
		expHand:    []Card{NewCardFromString(`AH`), NewCardFromString(`2H`), NewCardFromString(`3H`), NewCardFromString(`4H`)},
	}, {
		msg:        `last 4 cards in the deck`,
		inputInt32: 0x30313233,
		expHand:    []Card{NewCardFromString(`10H`), NewCardFromString(`JH`), NewCardFromString(`QH`), NewCardFromString(`KH`)},
	}, {
		msg:        `2S,3s,4s,5s`,
		inputInt32: 0x01020304,
		expHand:    []Card{NewCardFromString(`2s`), NewCardFromString(`3s`), NewCardFromString(`4s`), NewCardFromString(`5s`)},
	}}

	for _, tc := range testCases {
		assert.Equal(t, tc.expHand, Int32ToHand(tc.inputInt32), tc.msg)
	}
}
