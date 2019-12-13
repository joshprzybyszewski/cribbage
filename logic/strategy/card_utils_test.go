package strategy

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/joshprzybyszewski/cribbage/model"
)

func isIn(c model.Card, h []model.Card) bool {
	for _, thisCard := range h {
		if c == thisCard {
			return true
		}
	}
	return false
}

func validateHand(origHand, thisHand []model.Card) bool {
	for _, c := range thisHand {
		if isIn(c, origHand) {
			return true
		}
	}
	return false
}

func TestChooseFrom(t *testing.T) {
	tests := []struct {
		desc   string
		hand   []model.Card
		nCards int
	}{{
		desc: ``,
		hand: []model.Card{
			model.NewCardFromString(`ac`),
			model.NewCardFromString(`2c`),
			model.NewCardFromString(`3c`),
			model.NewCardFromString(`4c`),
			model.NewCardFromString(`5c`),
			model.NewCardFromString(`6c`),
		},
	}}
	for _, tc := range tests {
		all := chooseFrom(tc.nCards, tc.hand)
		for _, h := range all {
			ok := validateHand(tc.hand, h)
			assert.True(t, ok)
		}
	}
}
