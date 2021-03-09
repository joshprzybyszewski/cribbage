// +build !prod

package network

import "github.com/joshprzybyszewski/cribbage/model"

func ModelCardsFromStrings(cs ...string) []model.Card {
	hand := make([]model.Card, len(cs))
	for i, c := range cs {
		hand[i] = model.NewCardFromString(c)
	}
	return hand
}
