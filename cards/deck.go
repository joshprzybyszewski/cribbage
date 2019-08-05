package cards

import (
	"crypto/rand"
	"math/big"
)

type Deck struct {
	cards    [52]Card
	numDealt int
}

func NewDeck() *Deck {
	var cards [52]Card
	i := 0
	for s := Spades; s <= Hearts; s++ {
		for v := 1; v <= 13; v++ {
			cards[i] = NewCard(s, v)
			i++
		}
	}

	return &Deck{
		cards:    cards,
		numDealt: 0,
	}
}

func (d *Deck) Deal() Card {
	lastValidCard := int64(51 - d.numDealt)
	randBigInt, err := rand.Int(rand.Reader, big.NewInt(lastValidCard))
	if err != nil {
		// rand.Int should never fail
		panic(err)
	}
	randomIndex := randBigInt.Int64()

	d.cards[lastValidCard], d.cards[randomIndex] = d.cards[randomIndex], d.cards[lastValidCard]

	d.numDealt++

	return d.cards[lastValidCard]
}

func (d *Deck) Shuffle() {
	for i := 0; i < 52; i++ {
		d.Deal()
	}

	d.numDealt = 0
}
