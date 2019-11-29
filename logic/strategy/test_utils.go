//+build !prod

package strategy

import (
	"github.com/joshprzybyszewski/cribbage/model"
)

func strToCards(s []string) []model.Card {
	c := make([]model.Card, len(s))
	for i, str := range s {
		c[i] = model.NewCardFromString(str)
	}
	return c
}

func strToPeggedCards(s []string) []model.PeggedCard {
	c := make([]model.PeggedCard, len(s))
	for i, str := range s {
		c[i] = model.NewPeggedCard(model.InvalidPlayerID, model.NewCardFromString(str), 0)
	}
	return c
}

func containsCard(cs []string, c model.Card) bool {
	for _, cstr := range cs {
		if model.NewCardFromString(cstr).String() == c.String() {
			return true
		}
	}
	return false
}
