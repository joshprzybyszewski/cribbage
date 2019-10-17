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
