package round

import (
	"errors"

	"github.com/joshprzybyszewski/cribbage/cards"
)

type Player struct {
	cards []cards.Card
	name  string

	TotalScore int
}

func NewPlayer(name string) *Player {
	return &Player{
		name:       name,
		cards:      make([]cards.Card, 0, 6),
		TotalScore: 0,
	}
}

func (p *Player) AcceptCard(c cards.Card) error {
	if len(cards) >= 6 {
		return errors.New(`cannot accept new card`)
	}

	p.cards = append(p.cards, c)
	return nil
}
