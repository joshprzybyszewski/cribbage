package strategy

import (
	"github.com/joshprzybyszewski/cribbage/cards"
)

// AvoidCribFifteens tries to return a set of cards which does not add up to 15
func AvoidCribFifteens(desired int, hand []cards.Card) []cards.Card {
	objFunc := func(c1, c2 cards.Card) bool {
		return c1.PegValue()+c2.PegValue() != 15
	}
	return determineCribCards(desired, hand, objFunc)
}

// GiveCribFifteens tries to return a set of cards which adds up to 15
func GiveCribFifteens(desired int, hand []cards.Card) []cards.Card {
	objFunc := func(c1, c2 cards.Card) bool {
		return c1.PegValue()+c2.PegValue() == 15
	}
	return determineCribCards(desired, hand, objFunc)
}

// AvoidCribPairs tries to return a set of cards which does not make a pair (unequal value)
func AvoidCribPairs(desired int, hand []cards.Card) []cards.Card {
	objFunc := func(c1, c2 cards.Card) bool {
		return c1.Value != c2.Value
	}
	return determineCribCards(desired, hand, objFunc)
}

// GiveCribPairs tries to return a set of cards that makes a pair (equal value)
func GiveCribPairs(desired int, hand []cards.Card) []cards.Card {
	objFunc := func(c1, c2 cards.Card) bool {
		return c1.Value == c2.Value
	}
	return determineCribCards(desired, hand, objFunc)
}

func determineCribCards(desired int, hand []cards.Card, objectiveFunc func(c1, c2 cards.Card) bool) []cards.Card {
	// Currently this function uses a very piecewise solution... but it passes the tests :)
	cribCards := make([]cards.Card, 0, desired)
	if desired == 1 {
		cribCards = append(cribCards, hand[0])
		return cribCards
	}
	for i := 0; i < len(hand)-1; i++ {
		c1 := hand[i]
		otherCards := hand[i+1:]
		for _, c2 := range otherCards {
			if objectiveFunc(c1, c2) {
				return append(cribCards, c1, c2)
			}
		}
	}
	return hand[:2]
}
