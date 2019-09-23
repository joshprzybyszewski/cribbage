package game

import (
	"github.com/joshprzybyszewski/cribbage/cards"
)

var _ Player = (*humanPlayer)(nil)

type humanPlayer struct {
	*player
}

func NewHumanPlayer(name string, color PegColor) Player {
	return &humanPlayer{
		player: newPlayer(name, color),
	}
}

func (p *humanPlayer) Shuffle() {
	// TODO ask the user how many times to shuffle
}

func (p *humanPlayer) AddToCrib() []cards.Card {
	// TODO ask the user which cards to add to the crib
	// return those
	// then update the users hand to not have them
	return nil
}

func (p *humanPlayer) Cut() float64 {
	// TODO Ask the user how far down the deck they wanna cut
	return 0
}

func (p *humanPlayer) Peg(maxVal int) (cards.Card, bool) {
	// TODO ask the user which of their cards they would like to peg
	// TODO valideate they can peg that
	return cards.Card{}, false
}
