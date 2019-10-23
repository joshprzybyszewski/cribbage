package strategy

import (
	"github.com/joshprzybyszewski/cribbage/cards"
)

func PegToFifteen(hand, prevPegs []cards.Card, curPeg int) (_ cards.Card, sayGo bool) {
	return cards.Card{}, true
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
