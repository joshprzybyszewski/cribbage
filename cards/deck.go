package cards

import (
	"crypto/rand"
	"math/big"
)

const NumCardsPerDeck = 52

type Deck struct {
	cards    [52]Card
	numDealt int
}

func NewDeck() *Deck {
	var cards [52]Card

	for i := 0; i < NumCardsPerDeck; i++ {
		cards[i] = NewCardFromNumber(i)
	}

	return &Deck{
		cards:    cards,
		numDealt: 0,
	}
}

func (d *Deck) Deal() Card {
	lastValidCard := int64(51 - d.numDealt)
	if lastValidCard > 0 {
		randBigInt, err := rand.Int(rand.Reader, big.NewInt(lastValidCard))
		if err != nil {
			// rand.Int should never fail
			panic(err)
		}
		randomIndex := randBigInt.Int64()
		d.cards[lastValidCard], d.cards[randomIndex] = d.cards[randomIndex], d.cards[lastValidCard]
	}

	d.numDealt++
	if d.numDealt > 52 {
		println(`bad time`)
	}

	return d.cards[lastValidCard]
}

func (d *Deck) Shuffle() {
	d.numDealt = 0
	
	for i := 0; i < 52; i++ {
		d.Deal()
	}

	d.numDealt = 0
}

func (d *Deck) CutDeck(p float64) Card {
	lastValidCard := int64(51 - d.numDealt)
	cutCard := int(float64(lastValidCard) * p)

	// say that all of the cards are dealt because once it's cut, we can't reuse it
	d.numDealt = 51

	return d.cards[cutCard]
}
