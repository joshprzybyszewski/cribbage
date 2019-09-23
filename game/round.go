package game

import (
	"errors"

	"github.com/joshprzybyszewski/cribbage/cards"
)

type Round struct {
	CurrentStage RoundStage

	// the cards that have been placed in the crib for this round
	cribCards []cards.Card

	// the ordered list of cards which is added to as players play a card during pegging
	peggedCards []cards.Card

	// the number we are currently at in pegging
	currentPeg int
}

func NewTwoPlayerRound() *Round {
	return newRound(nil, 2)
}

func NewThreePlayerRound(cribCard cards.Card) *Round {
	cc := make([]cards.Card, 0, 4)
	cc = append(cc, cribCard)

	return newRound(cc, 3)
}

func NewFourPlayerRound() *Round {
	return newRound(nil, 4)
}

func newRound(cribCards []cards.Card, numPlayers int) *Round {
	cc := make([]cards.Card, 0, 4)
	if cribCards != nil {
		cc = append(cc, cribCards...)
	}

	return &Round{
		CurrentStage: Deal,
		cribCards:    cc,
		peggedCards:  make([]cards.Card, 0, 4*numPlayers),
		currentPeg:   0,
	}
}
