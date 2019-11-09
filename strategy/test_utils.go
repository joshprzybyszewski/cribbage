//+build !prod

package strategy

import (
	"github.com/joshprzybyszewski/cribbage/cards"
)

func strToCards(s []string) []cards.Card {
	c := make([]cards.Card, len(s))
	for i, str := range s {
		c[i] = cards.NewCardFromString(str)
	}
	return c
}

func containsCard(cs []string, c cards.Card) bool {
	for _, cstr := range cs {
		if cards.NewCardFromString(cstr).String() == c.String() {
			return true
		}
	}
	return false
}
