package network

import "github.com/joshprzybyszewski/cribbage/model"

func newCardFromModel(c model.Card) Card {
	return Card{
		Suit:  c.Suit.String(),
		Value: c.Value,
		Name:  c.String(),
	}
}
