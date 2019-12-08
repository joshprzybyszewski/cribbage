package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	fakeCard := Card{}

	dealtCards := map[string]struct{}{}
	for i := 0; i < 55; i++ {
		c := d.Deal()
		if _, ok := dealtCards[c.String()]; ok && c != fakeCard {
			t.Fatalf("Should not have dealt the same card twice: %s", c.String())
		}
		dealtCards[c.String()] = struct{}{}
		if i >= NumCardsPerDeck {
			assert.Equal(t, Card{}, c)
		} else {
			assert.NotEqual(t, Card{}, c)
		}
	}
}

func TestDeckCutting(t *testing.T) {
	d := NewDeck()

	cutCard1, err := d.CutDeck(0.5)
	require.NoError(t, err)
	cutCard2, err := d.CutDeck(0.5)
	require.NoError(t, err)
	assert.Equal(t, cutCard1, cutCard2)
	cutCard3, err := d.CutDeck(0.5)
	require.NoError(t, err)
	assert.Equal(t, cutCard1, cutCard3)
	assert.Equal(t, cutCard2, cutCard3)
	cutCard4, err := d.CutDeck(0.75)
	require.NoError(t, err)
	assert.NotEqual(t, cutCard1, cutCard4)
}

func TestDeckCuttingAvoidsDealtCards(t *testing.T) {
	d := NewDeck()
	dealtCard1 := d.Deal()
	dealtCard2 := d.Deal()
	for i := 0; i <= 104; i++ {
		p := float64(i) / 104.0
		cutCard, err := d.CutDeck(p)
		require.NoError(t, err)
		assert.NotEqual(t, dealtCard1, cutCard)
		assert.NotEqual(t, dealtCard2, cutCard)
	}
}
