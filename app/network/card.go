package network

import "github.com/joshprzybyszewski/cribbage/model"

var (
	invalidCard Card = Card{
		Suit:  ``,
		Value: -1,
		Name:  `unknown`,
	}
)

type Card struct {
	Suit  string `json:"suit"`
	Value int    `json:"value"`
	Name  string `json:"name"`
}

func convertToCards(mCards []model.Card) []Card {
	if mCards == nil {
		return nil
	}
	cards := make([]Card, len(mCards))
	for i, c := range mCards {
		cards[i] = convertToCard(c)
	}
	return cards
}

func convertFromCards(cs []Card) []model.Card {
	if cs == nil {
		return nil
	}
	mcs := make([]model.Card, len(cs))
	for i, c := range cs {
		mcs[i] = convertFromCard(c)
	}
	return mcs
}

func convertToCard(c model.Card) Card {
	return Card{
		Suit:  c.Suit.String(),
		Value: c.Value,
		Name:  c.String(),
	}
}

func convertFromCard(c Card) model.Card {
	return model.NewCardFromString(c.Name)
}

type PeggedCard struct {
	Card   Card           `json:"card"`
	Player model.PlayerID `json:"player"`
}

func convertToPeggedCards(mPeggedCards []model.PeggedCard) []PeggedCard {
	if mPeggedCards == nil {
		return nil
	}
	cards := make([]PeggedCard, len(mPeggedCards))
	for i, pc := range mPeggedCards {
		cards[i].Card = convertToCard(pc.Card)
		cards[i].Player = pc.PlayerID
	}
	return cards
}

func convertFromPeggedCards(pcs []PeggedCard) []model.PeggedCard {
	if pcs == nil {
		return nil
	}
	mpcs := make([]model.PeggedCard, len(pcs))
	for i, pc := range pcs {
		mpcs[i].Card = convertFromCard(pc.Card)
		mpcs[i].Action = i
		mpcs[i].PlayerID = pc.Player
	}
	return mpcs
}
