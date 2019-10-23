package strategy

import (
	"github.com/joshprzybyszewski/cribbage/cards"
)

// PegToFifteen returns a card that yields a fifteen if it can
func PegToFifteen(hand, prevPegs []cards.Card, curPeg int) (_ cards.Card, sayGo bool) {
	ct := 0
	for _, c := range hand {
		pegVal := c.PegValue() + curPeg
		switch {
		case pegVal == 15:
			return c, false
		case pegVal > 31:
			ct++
		}
	}
	if ct == len(hand) {
		return cards.Card{}, true
	}
	return hand[0], false
}

func PegToThirtyOne(hand, prevPegs []cards.Card, curPeg int) (_ cards.Card, sayGo bool) {
	return cards.Card{}, true
}

func PegToPair(hand, prevPegs []cards.Card, curPeg int) (_ cards.Card, sayGo bool) {
	return cards.Card{}, true
}

func PegToRun(hand, prevPegs []cards.Card, curPeg int) (_ cards.Card, sayGo bool) {
	return cards.Card{}, true
}
