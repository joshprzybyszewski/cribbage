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

func Test_newDeckWithDealt(t *testing.T) {
	already := map[Card]struct{}{
		Card{
			Suit:  Spades,
			Value: 2,
		}: struct{}{},
	}
	dealtDeck := newDeckWithDealt(already)

	d, ok := dealtDeck.(*deck)
	require.True(t, ok)
	assert.Equal(t, len(already), d.numDealt)
	assert.Equal(t, Card{
		Suit:  Spades,
		Value: 2,
	}, d.cards[51])

	already[Card{
		Suit:  Hearts,
		Value: 6,
	}] = struct{}{}
	dealtDeck = newDeckWithDealt(already)

	d, ok = dealtDeck.(*deck)
	require.True(t, ok)
	assert.Equal(t, len(already), d.numDealt)
	assert.ElementsMatch(t, []Card{{
		Suit:  Spades,
		Value: 2,
	}, {
		Suit:  Hearts,
		Value: 6,
	}}, d.cards[50:])

	already[Card{
		Suit:  Clubs,
		Value: 5,
	}] = struct{}{}
	dealtDeck = newDeckWithDealt(already)

	d, ok = dealtDeck.(*deck)
	require.True(t, ok)
	assert.Equal(t, len(already), d.numDealt)
	assert.ElementsMatch(t, []Card{{
		Suit:  Spades,
		Value: 2,
	}, {
		Suit:  Hearts,
		Value: 6,
	}, {
		Suit:  Clubs,
		Value: 5,
	}}, d.cards[49:])
}

func TestGetDeck(t *testing.T) {
	testCases := []struct {
		msg         string
		game        Game
		expNumDealt int
	}{{
		msg:         `empty game`,
		game:        Game{},
		expNumDealt: 0,
	}, {
		msg: `with hands`,
		game: Game{
			Hands: map[PlayerID][]Card{
				PlayerID(`alice`): []Card{
					NewCardFromString(`1s`),
					NewCardFromString(`2s`),
					NewCardFromString(`3s`),
					NewCardFromString(`4s`),
					NewCardFromString(`5s`),
					NewCardFromString(`6s`),
				},
				PlayerID(`bob`): []Card{
					NewCardFromString(`1c`),
					NewCardFromString(`2c`),
					NewCardFromString(`3c`),
					NewCardFromString(`4c`),
					NewCardFromString(`5c`),
					NewCardFromString(`6c`),
				},
			},
		},
		expNumDealt: 12,
	}, {
		msg: `with hands and crib`,
		game: Game{
			Hands: map[PlayerID][]Card{
				PlayerID(`alice`): []Card{
					NewCardFromString(`1s`),
					NewCardFromString(`2s`),
					NewCardFromString(`3s`),
					NewCardFromString(`4s`),
				},
				PlayerID(`bob`): []Card{
					NewCardFromString(`1c`),
					NewCardFromString(`2c`),
					NewCardFromString(`3c`),
					NewCardFromString(`4c`),
				},
			},
			Crib: []Card{
				NewCardFromString(`js`),
				NewCardFromString(`jc`),
				NewCardFromString(`jd`),
				NewCardFromString(`jh`),
			},
		},
		expNumDealt: 12,
	}}

	for _, tc := range testCases {
		g := tc.game
		gDeck, err := g.GetDeck()
		require.NoError(t, err, tc.msg)

		d, ok := gDeck.(*deck)
		require.True(t, ok, tc.msg)

		d1 := *d

		assert.Equal(t, tc.expNumDealt, d1.numDealt, tc.msg)

		gDeck, err = g.GetDeck()
		require.NoError(t, err, tc.msg)

		d, ok = gDeck.(*deck)
		require.True(t, ok, tc.msg)

		d2 := *d

		assert.Equal(t, d1.numDealt, d2.numDealt, tc.msg)
		assert.ElementsMatch(t, d1.cards[:len(d1.cards)-d1.numDealt], d2.cards[:len(d1.cards)-d1.numDealt], tc.msg)
	}

}
