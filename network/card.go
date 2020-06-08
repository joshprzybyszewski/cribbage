package network

import "github.com/joshprzybyszewski/cribbage/model"

type Card struct {
	Suit  string `json:"suit"`
	Value int    `json:"value"`
	Name  string `json:"name"`
}

func convertToCards(mCards []model.Card) []Card {
	cards := make([]Card, len(mCards))
	for i, c := range mCards {
		cards[i] = convertToCard(c)
	}
	return cards
}

func convertToCard(c model.Card) Card {
	return Card{
		Suit:  c.Suit.String(),
		Value: c.Value,
		Name:  c.String(),
	}
}

type PeggedCard struct {
	Card   Card           `json:"card"`
	Player model.PlayerID `json:"player"`
}

func convertPeggedCards(mPeggedCards []model.PeggedCard) []PeggedCard {
	cards := make([]PeggedCard, len(mPeggedCards))
	for i, pc := range mPeggedCards {
		cards[i].Card = convertToCard(pc.Card)
		cards[i].Player = pc.PlayerID
	}
	return cards
}
