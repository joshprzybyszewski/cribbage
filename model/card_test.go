package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJackValue(t *testing.T) {
	assert.Equal(t, JackValue, NewCardFromString(`jh`).Value)
}

func TestNewCardFromString(t *testing.T) {
	testCases := []struct {
		desc    string
		input   string
		expCard Card
	}{{
		desc:  `ace of hearts`,
		input: `AH`,
		expCard: Card{
			Suit:  Hearts,
			Value: 1,
		},
	}, {
		desc:  `two of diamonds`,
		input: `2D`,
		expCard: Card{
			Suit:  Diamonds,
			Value: 2,
		},
	}, {
		desc:  `10 of spades`,
		input: `10S`,
		expCard: Card{
			Suit:  Spades,
			Value: 10,
		},
	}, {
		desc:  `Jack of Clubs`,
		input: `JC`,
		expCard: Card{
			Suit:  Clubs,
			Value: 11,
		},
	}, {
		desc:  `Queen of Hearts`,
		input: `QH`,
		expCard: Card{
			Suit:  Hearts,
			Value: 12,
		},
	}, {
		desc:  `King of Diamonds`,
		input: `KD`,
		expCard: Card{
			Suit:  Diamonds,
			Value: 13,
		},
	}, {
		desc:  `10 of spades`,
		input: `10s`,
		expCard: Card{
			Suit:  Spades,
			Value: 10,
		},
	}, {
		desc:  `Jack of Clubs`,
		input: `11c`,
		expCard: Card{
			Suit:  Clubs,
			Value: 11,
		},
	}}

	for _, tc := range testCases {
		assert.Equal(t, tc.expCard, NewCardFromString(tc.input), tc.desc)
	}
}

func TestNewCardFromStringWithWeirdInput(t *testing.T) {
	// We don't support emojis
	testCases := []struct {
		desc  string
		input string
	}{{
		desc:  `ace of hearts`,
		input: `A♡`,
	}, {
		desc:  `two of diamonds`,
		input: `2♢`,
	}, {
		desc:  `Queen of Hearts`,
		input: `12♥︎`,
	}, {
		desc:  `King of Diamonds`,
		input: `13♦`,
	}}

	for _, tc := range testCases {
		assert.Equal(t, Card{}, NewCardFromString(tc.input), tc.desc)
	}
}

func TestPegValue(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		expValue int
	}{{
		desc:     `ace of hearts`,
		input:    `AH`,
		expValue: 1,
	}, {
		desc:     `two of diamonds`,
		input:    `2D`,
		expValue: 2,
	}, {
		desc:     `10 of spades`,
		input:    `10s`,
		expValue: 10,
	}, {
		desc:     `Jack of Clubs`,
		input:    `11c`,
		expValue: 10,
	}, {
		desc:     `Queen of Hearts`,
		input:    `12h`,
		expValue: 10,
	}, {
		desc:     `King of Diamonds`,
		input:    `13d`,
		expValue: 10,
	}}

	for _, tc := range testCases {
		assert.Equal(t, tc.expValue, NewCardFromString(tc.input).PegValue(), tc.desc)
	}
}

func TestNewCardFromTinyInt(t *testing.T) {
	for ti := int8(0); ti < 52; ti++ {
		c, err := NewCardFromTinyInt(ti)
		require.NoError(t, err, `should not error for value %d`, ti)
		outputTI := c.ToTinyInt()
		assert.Equal(t, ti, outputTI)
	}

	ti := int8(52)
	c, err := NewCardFromTinyInt(ti)
	assert.Error(t, err, `should error for value %d`, ti)
	outputTI := c.ToTinyInt()
	assert.NotEqual(t, ti, outputTI)
	assert.Equal(t, int8(-1), outputTI)

	ti = int8(-2)
	c, err = NewCardFromTinyInt(ti)
	assert.Error(t, err, `should error for value %d`, ti)
	outputTI = c.ToTinyInt()
	assert.NotEqual(t, ti, outputTI)
	assert.Equal(t, int8(-1), outputTI)
}
