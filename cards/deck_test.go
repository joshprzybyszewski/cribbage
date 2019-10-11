package cards

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeckBasics(t *testing.T) {
	assert.Equal(t, 52, NumCardsPerDeck)
}

func TestDeckShuffling(t *testing.T) {
	d := NewDeck()

	for i := 0; i < 100; i++ {
		d.Shuffle()
	}
}

func TestDeckDealing(t *testing.T) {
	d := NewDeck()

	dealtCards := map[string]struct{}{}
	for i := 0; i < 55; i++ {
		c := d.Deal()
		if _, ok := dealtCards[c.String()]; ok && c != Card{} {
			t.Fatalf("Should not have dealt the same card twice: %s", c.String())
		}
		dealtCards[c.String()] = struct{}{}
		if i < NumCardsPerDeck {
			assert.Equal(t, Card{}, c)
		} else {
			assert.NotEqual(t, Card{}, c)
		}
	}
}

func TestDeckCutting(t *testing.T) {
	d := NewDeck()

	cutCard1 := d.CutDeck(0.5)
	cutCard2 := d.CutDeck(0.5)
	assert.NotEqual(t, cutCard1, cutCard2)
	cutCard3 := d.CutDeck(0.5)
	assert.NotEqual(t, cutCard1, cutCard3)
	assert.Equal(t, cutCard2, cutCard3, `after the first cut, nothing is sane`)

	d.Shuffle()
	cutCard4 := d.CutDeck(0.5)
	assert.NotEqual(t, cutCard1, cutCard4, `this has a _very low_ chance of being equal`)
}
