package strategy

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/joshprzybyszewski/cribbage/cards"
)

func TestPegToFifteen(t *testing.T) {
	testCases := []struct {
		msg           string
		inputHand     []string
		inputPrevPegs []string
		inputCurPeg   int
		expCard       cards.Card
		expSayGo      bool
	}{{
		msg:           `obvious`,
		inputHand:     []string{`3c`, `4d`, `5s`, `6h`},
		inputPrevPegs: []string{`10s`},
		inputCurPeg:   10,
		expCard:       cards.NewCardFromString(`5s`),
		expSayGo:      false,
	}}

	for _, tc := range testCases {
		actCard, actSayGo := PegToFifteen(strToCards(tc.inputHand), strToCards(tc.inputPrevPegs), tc.inputCurPeg)
		assert.Equal(t, tc.expSayGo, actSayGo, tc.msg)
		assert.Equal(t, tc.expCard, actCard, tc.msg)
	}
}

func TestPegToThirtyOne(t *testing.T) {
	testCases := []struct {
		msg           string
		inputHand     []string
		inputPrevPegs []string
		inputCurPeg   int
		expCard       cards.Card
		expSayGo      bool
	}{{
		msg:           `obvious`,
		inputHand:     []string{`ac`, `2h`, `4s`},
		inputPrevPegs: []string{`10s`, `10h`, `10d`},
		inputCurPeg:   30,
		expCard:       cards.NewCardFromString(`ac`),
		expSayGo:      false,
	}}

	for _, tc := range testCases {
		actCard, actSayGo := PegToThirtyOne(strToCards(tc.inputHand), strToCards(tc.inputPrevPegs), tc.inputCurPeg)
		assert.Equal(t, tc.expSayGo, actSayGo, tc.msg)
		assert.Equal(t, tc.expCard, actCard, tc.msg)
	}
}

func TestPegToPair(t *testing.T) {
	testCases := []struct {
		msg           string
		inputHand     []string
		inputPrevPegs []string
		inputCurPeg   int
		expCard       cards.Card
		expSayGo      bool
	}{{
		msg:           `obvious`,
		inputHand:     []string{`3c`, `10s`, `6d`, `2h`},
		inputPrevPegs: []string{`3d`},
		inputCurPeg:   3,
		expCard:       cards.NewCardFromString(`3c`),
		expSayGo:      false,
	}}

	for _, tc := range testCases {
		actCard, actSayGo := PegToPair(strToCards(tc.inputHand), strToCards(tc.inputPrevPegs), tc.inputCurPeg)
		assert.Equal(t, tc.expSayGo, actSayGo, tc.msg)
		assert.Equal(t, tc.expCard, actCard, tc.msg)
	}
}

func TestPegToRun(t *testing.T) {
	testCases := []struct {
		msg           string
		inputHand     []string
		inputPrevPegs []string
		inputCurPeg   int
		expCard       cards.Card
		expSayGo      bool
	}{{
		msg:           `obvious`,
		inputHand:     []string{`2d`, `7h`, `jc`},
		inputPrevPegs: []string{`8h`, `6d`},
		inputCurPeg:   14,
		expCard:       cards.NewCardFromString(`7h`),
		expSayGo:      false,
	}}

	for _, tc := range testCases {
		actCard, actSayGo := PegToRun(strToCards(tc.inputHand), strToCards(tc.inputPrevPegs), tc.inputCurPeg)
		assert.Equal(t, tc.expSayGo, actSayGo, tc.msg)
		assert.Equal(t, tc.expCard, actCard, tc.msg)
	}
}
