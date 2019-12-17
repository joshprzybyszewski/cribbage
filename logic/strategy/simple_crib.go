package strategy

import (
	"github.com/joshprzybyszewski/cribbage/model"
)

// AvoidCribFifteens tries to return a set of cards which does not add up to 15
func AvoidCribFifteens(desired int, hand []model.Card) []model.Card {
	objFunc := func(c1, c2 model.Card) bool {
		return c1.PegValue()+c2.PegValue() != 15
	}
	return determineCribCards(desired, hand, objFunc)
}

// GiveCribFifteens tries to return a set of cards which adds up to 15
func GiveCribFifteens(desired int, hand []model.Card) []model.Card {
	objFunc := func(c1, c2 model.Card) bool {
		return c1.PegValue()+c2.PegValue() == 15
	}
	return determineCribCards(desired, hand, objFunc)
}

// AvoidCribPairs tries to return a set of cards which does not make a pair (unequal value)
func AvoidCribPairs(desired int, hand []model.Card) []model.Card {
	objFunc := func(c1, c2 model.Card) bool {
		return c1.Value != c2.Value
	}
	return determineCribCards(desired, hand, objFunc)
}

// GiveCribPairs tries to return a set of cards that makes a pair (equal value)
func GiveCribPairs(desired int, hand []model.Card) []model.Card {
	objFunc := func(c1, c2 model.Card) bool {
		return c1.Value == c2.Value
	}
	return determineCribCards(desired, hand, objFunc)
}

func determineCribCards(desired int, hand []model.Card, objectiveFunc func(c1, c2 model.Card) bool) []model.Card {
	if desired == 1 {
		return []model.Card{hand[0]}
	}
	for i := 0; i < len(hand)-1; i++ {
		c1 := hand[i]
		otherCards := hand[i+1:]
		for _, c2 := range otherCards {
			if objectiveFunc(c1, c2) {
				return []model.Card{c1, c2}
			}
		}
	}
	return hand[:2]
}
