package game

import (
	"errors"

	"github.com/joshprzybyszewski/cribbage/cards"
	"github.com/joshprzybyszewski/cribbage/scorer"
)

type Player interface {
	Name() string
	Color() PegColor

	TakeDeck(d *cards.Deck)
	IsDealer() bool
	Shuffle()

	DealCard() (cards.Card, error)
	AcceptCard(cards.Card) error
	NeedsCard() bool

	AddToCrib() []cards.Card
	AcceptCrib([]cards.Card) error

	Cut() float64

	Peg(maxVal int) (played cards.Card, sayGo, canPlay bool)

	HandScore(leadCard cards.Card) int

	CribScore(leadCard cards.Card) (int, error)

	PassDeck() error
}

type player struct {
	name  string
	color PegColor

	deck *cards.Deck

	hand   []cards.Card
	pegged map[cards.Card]struct{}
	crib   []cards.Card
}

func newPlayer(name string, color PegColor) *player {
	return &player{
		name:  name,
		color: color,
		deck:  nil,
		hand:  make([]cards.Card, 0, 6),
		crib:  make([]cards.Card, 0, 4),
	}
}

func (p *player) Name() string {
	return p.name
}

func (p *player) Color() PegColor {
	return p.color
}

func (p *player) TakeDeck(d *cards.Deck) {
	p.deck = d
}

func (p *player) IsDealer() bool {
	return p.deck != nil
}

func (p *player) DealCard() (cards.Card, error) {
	if !p.IsDealer() {
		return cards.Card{}, errors.New(`cannot deal a card if not the dealer`)
	}

	return p.deck.Deal(), nil
}

func (p *player) AcceptCard(c cards.Card) error {
	if !p.NeedsCard() {
		return errors.New(`cannot accept new card`)
	}

	p.hand = append(p.hand, c)
	return nil
}

func (p *player) NeedsCard() bool {
	// TODO this is only currently setup for 2 player games
	return len(p.hand) < 6
}

func (p *player) AcceptCrib([]cards.Card) error {
	if !p.IsDealer() {
		return errors.New(`attempted to receive a crib when not the dealer`)
	}

	return nil
}

func (p *player) HandScore(leadCard cards.Card) int {
	return scorer.CribPoints(leadCard, p.hand)
}

func (p *player) CribScore(leadCard cards.Card) (int, error) {
	if !p.IsDealer() {
		return 0, errors.New(`Cannot score crib when not the dealer`)
	} else if len(p.crib) == 0 {
		return 0, errors.New(`expected to have crib, but missing!`)
	}

	return scorer.CribPoints(leadCard, p.crib), nil
}

func (p *player) PassDeck() error {
	if p.IsDealer() {
		return errors.New(`cannot pass the deck when not the dealer`)
	} else if p.deck == nil {
		return errors.New(`didn't have deck!`)
	}

	p.deck = nil
	p.crib = p.crib[:0]

	return nil
}
