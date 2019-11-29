package model

import (
	"github.com/joshprzybyszewski/cribbage/utils/rand"
)

type Deck interface {
	Deal() Card
	Shuffle()
	CutDeck(p float64) Card
}

type deck struct {
	cards    [52]Card
	numDealt int
}

func NewDeck() Deck {
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

func (d *deck) Deal() Card {
	lastValidCard := 51 - d.numDealt
	if lastValidCard > 0 {
		randomIndex := rand.Intn(lastValidCard)
		d.cards[lastValidCard], d.cards[randomIndex] = d.cards[randomIndex], d.cards[lastValidCard]
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

func (d *deck) CutDeck(p float64) Card {
	lastValidCard := int64(51 - d.numDealt)
	cutCard := int(float64(lastValidCard) * p)

	// say that all of the cards are dealt because once it's cut, we can't reuse it
	d.numDealt = 51

	return d.cards[cutCard]
}
