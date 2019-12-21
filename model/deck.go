package model

import (
	"errors"

	"github.com/joshprzybyszewski/cribbage/utils/rand"
)

type Deck interface {
	Deal() Card
	Shuffle()
	CutDeck(p float64) (Card, error)
}

type deck struct {
	cards    [52]Card
	numDealt int
}

func NewDeck() Deck {
	return newDeck()
}

func newDeck() *deck {
	var cards [52]Card

	for i := 0; i < NumCardsPerDeck; i++ {
		cards[i] = NewCardFromNumber(i)
	}

	d := deck{
		cards:    cards,
		numDealt: 0,
	}

	// Start the deck off in a random state by shuffling it a few times
	n := rand.Intn(10)
	for i := 0; i < n; i++ {
		d.Shuffle()
	}

	return &d
}

func newDeckWithDealt(dealt map[Card]struct{}) Deck {
	d := newDeck()
	if len(dealt) == 0 {
		return d
	}

	for i := 0; i < len(d.cards); i++ {
		c := d.cards[i]
		lastValidCard := 51 - d.numDealt
		if i >= lastValidCard {
			break
		}
		if _, ok := dealt[c]; ok {
			tmp := d.cards[lastValidCard]
			d.cards[lastValidCard] = d.cards[i]
			d.cards[i] = tmp
			d.numDealt++
			i--
		}
	}

	return d
}

func (d *deck) Deal() Card {
	lastValidCard := 51 - d.numDealt
	if lastValidCard > 0 {
		randomIndex := rand.Intn(lastValidCard)
		tmp := d.cards[lastValidCard]
		d.cards[lastValidCard] = d.cards[randomIndex]
		d.cards[randomIndex] = tmp
	}

	d.numDealt++
	if d.numDealt > 52 {
		// This is a :badtime:
		return Card{}
	}

	return d.cards[lastValidCard]
}

func (d *deck) Shuffle() {
	d.numDealt = 0

	for i := 0; i < 52; i++ {
		d.Deal()
	}

	d.numDealt = 0
}

func (d *deck) CutDeck(p float64) (Card, error) {
	if d.numDealt >= 52 {
		return Card{}, errors.New(`cannot cut deck with all cards dealt`)
	}

	lastValidCard := int64(51 - d.numDealt)
	cutCard := int(float64(lastValidCard) * p)

	return d.cards[cutCard], nil
}
