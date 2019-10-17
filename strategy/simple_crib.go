package strategy

import (
	"github.com/joshprzybyszewski/cribbage/cards"
)

// AvoidCribFifteens tries to return a set of cards which does not add up to 15
func AvoidCribFifteens(desired int, hand []cards.Card) []cards.Card {
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
			if c1.PegValue()+c2.PegValue() != 15 {
				return append(cribCards, c1, c2)
			}
		}
	}
	return append(cribCards, hand[0], hand[1])
}

func GiveCribFifteens(desired int, hand []cards.Card) []cards.Card {

	return nil
}

func prioritize(c cards.Card, value int) int {
	if c.Value == value {
		return 1
	}
	return 0
}

func AvoidCribPairs(desired int, hand []cards.Card) []cards.Card {

	return nil
}

func GiveCribPairs(desired int, hand []cards.Card) []cards.Card {

	return nil
}
