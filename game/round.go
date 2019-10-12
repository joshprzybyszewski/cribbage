package game

import (
	"errors"

	"github.com/joshprzybyszewski/cribbage/cards"
	"github.com/joshprzybyszewski/cribbage/pegging"
)

const (
	maxPeggingValue int = 31
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

func (r *Round) NextRound() error {
	if r.CurrentStage != Done {
		return errors.New(`cannot progress to next round when not done`)
	}

	r.CurrentStage = Deal
	r.cribCards = r.cribCards[:0]
	r.peggedCards = r.peggedCards[:0]
	r.currentPeg = 0
	return nil
}

func (r *Round) AcceptCribCards(c []cards.Card) error {
	if len(r.cribCards)+len(c) > 4 {
		return errors.New(`cannot accept cards -- crib would be too big`)
	}

	r.cribCards = append(r.cribCards, c...)
	return nil
}

func (r *Round) Crib() []cards.Card {
	return r.cribCards
}

func (r *Round) AcceptPegCard(c cards.Card) (int, error) {
	if r.currentPeg+c.PegValue() > maxPeggingValue {
		return 0, errors.New(`cannot peg past 31`)
	}

	pts, err := pegging.PointsForCard(r.peggedCards, c)
	if err != nil {
		return 0, err
	}

	r.peggedCards = append(r.peggedCards, c)
	r.currentPeg += c.PegValue()

	if 31 == r.currentPeg {
		r.currentPeg = 0
	}

	return pts, nil
}

func (r *Round) GoAround() {
	r.currentPeg = 0
}

func (r *Round) PrevPeggedCards() []cards.Card {
	return r.peggedCards
}

func (r *Round) CurrentPeg() int {
	return r.currentPeg
}
