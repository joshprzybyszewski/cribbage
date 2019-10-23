package strategy

import (
	"github.com/joshprzybyszewski/cribbage/cards"
)

// PegToFifteen returns a card that yields a fifteen if it can
func PegToFifteen(hand, prevPegs []cards.Card, curPeg int) (_ cards.Card, sayGo bool) {
	return pegToTarget(hand, prevPegs, curPeg, 15)
}

// PegToThirtyOne returns a card that yields 31 if it can
func PegToThirtyOne(hand, prevPegs []cards.Card, curPeg int) (_ cards.Card, sayGo bool) {
	return pegToTarget(hand, prevPegs, curPeg, 31)
}

func pegToTarget(hand, prevPegs []cards.Card, curPeg, target int) (_ cards.Card, sayGo bool) {
	ct := 0
	for _, c := range hand {
		pegVal := c.PegValue() + curPeg
		switch {
		case pegVal == target:
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

// PegToPair returns a card from the hand iff that card makes a pair and does not push the count over 31
func PegToPair(hand, prevPegs []cards.Card, curPeg int) (_ cards.Card, sayGo bool) {
	ct := 0
	lastCard := []rune(prevPegs[len(prevPegs)-1].String())
	for _, c := range hand {
		handCard := []rune(c.String())
		if curPeg+c.PegValue() > 31 {
			ct++
		} else if handCard[0] == lastCard[0] {
			return c, false
		}
	}
	if ct == len(hand) {
		return cards.Card{}, true
	}
	return hand[0], false
}

func PegToRun(hand, prevPegs []cards.Card, curPeg int) (_ cards.Card, sayGo bool) {
	return cards.Card{}, true
}
