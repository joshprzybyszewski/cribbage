package model

import (
	"crypto/rand"
	"fmt"
	"math/big"
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

	return &deck{
		cards:    cards,
		numDealt: 0,
	}
}

func (d *deck) Deal() Card {
	lastValidCard := int64(51 - d.numDealt)
	if lastValidCard > 0 {
		randBigInt, err := rand.Int(rand.Reader, big.NewInt(lastValidCard))
		if err != nil {
			// rand.Int should never fail
			fmt.Printf("Deal got error: %+v\n", err)
			return Card{}
		}
		randomIndex := randBigInt.Int64()
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
