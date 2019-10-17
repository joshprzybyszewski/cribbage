package strategy

import (
	"sort"

	"github.com/joshprzybyszewski/cribbage/cards"
)

// AvoidCribFifteens gg golint
func AvoidCribFifteens(desired int, hand []cards.Card) []cards.Card {
	cribCards := make([]cards.Card, 0, desired)
	sort.Slice(hand, func(i, j int) bool {
		return prioritize(hand[i], 5) < prioritize(hand[j], 5)
	})
	for i := 0; i < desired; i++ {
		cribCards = append(cribCards, hand[i])
	}
	// Sorting because the test is expecting the cards in a certain order...
	sort.Slice(cribCards, func(i, j int) bool {
		return cribCards[i].Value < cribCards[j].Value
	})
	return cribCards
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
