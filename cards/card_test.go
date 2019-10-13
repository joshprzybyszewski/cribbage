package cards

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCardFromString(t *testing.T) {
	testCases := []struct {
		desc    string
		input   string
		expCard Card
	}{{
		desc:  `ace of hearts`,
		input: `AH`,
		expCard: Card{
			Suit:      Hearts,
			Value:     1,
			deckValue: 39,
		},
	}, {
		desc:  `two of diamonds`,
		input: `2D`,
		expCard: Card{
			Suit:      Diamonds,
			Value:     2,
			deckValue: 27,
		},
	}, {
		desc:  `10 of spades`,
		input: `10S`,
		expCard: Card{
			Suit:      Spades,
			Value:     10,
			deckValue: 9,
		},
	}, {
		desc:  `Jack of Clubs`,
		input: `JC`,
		expCard: Card{
			Suit:      Clubs,
			Value:     11,
			deckValue: 23,
		},
	}, {
		desc:  `Queen of Hearts`,
		input: `QH`,
		expCard: Card{
			Suit:      Hearts,
			Value:     12,
			deckValue: 50,
		},
	}, {
		desc:  `King of Diamonds`,
		input: `KD`,
		expCard: Card{
			Suit:      Diamonds,
			Value:     13,
			deckValue: 38,
		},
	}}

	for _, tc := range testCases {
		assert.Equal(t, tc.expCard, NewCardFromString(tc.input), tc.desc)
	}
}

func TestNewCardFromStringWithWeirdInput(t *testing.T) {
	testCases := []struct {
		desc    string
		input   string
		expCard Card
	}{{
		desc:  `ace of hearts`,
		input: `A♡`,
		expCard: Card{
			Suit:      Hearts,
			Value:     1,
			deckValue: 39,
		},
	}, {
		desc:  `two of diamonds`,
		input: `2♢`,
		expCard: Card{
			Suit:      Diamonds,
			Value:     2,
			deckValue: 27,
		},
	}, {
		desc:  `10 of spades`,
		input: `10s`,
		expCard: Card{
			Suit:      Spades,
			Value:     10,
			deckValue: 9,
		},
	}, {
		desc:  `Jack of Clubs`,
		input: `11c`,
		expCard: Card{
			Suit:      Clubs,
			Value:     11,
			deckValue: 23,
		},
	}, {
		desc:  `Queen of Hearts`,
		input: `12♥︎`,
		expCard: Card{
			Suit:      Hearts,
			Value:     12,
			deckValue: 50,
		},
	}, {
		desc:  `King of Diamonds`,
		input: `13♦`,
		expCard: Card{
			Suit:      Diamonds,
			Value:     13,
			deckValue: 38,
		},
	}}

	for _, tc := range testCases {
		assert.Equal(t, tc.expCard, NewCardFromString(tc.input), tc.desc)
	}
}
