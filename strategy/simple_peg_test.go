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
	}, {
		msg:           `when can't, still returns a card`,
		inputHand:     []string{`10c`},
		inputPrevPegs: []string{`10s`, `3c`, `4d`, `5s`, `6h`, `3h`, `2h`},
		inputCurPeg:   2,
		expCard:       cards.NewCardFromString(`10c`),
		expSayGo:      false,
	}, {
		msg:           `says go when can't play`,
		inputHand:     []string{`3s`, `2h`, `4s`},
		inputPrevPegs: []string{`10s`, `10h`, `10d`},
		inputCurPeg:   30,
		expCard:       cards.Card{},
		expSayGo:      true,
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
	}, {
		msg:           `another less obvious`,
		inputHand:     []string{`7c`, `2h`, `4s`},
		inputPrevPegs: []string{`10s`, `10h`, `4d`},
		inputCurPeg:   24,
		expCard:       cards.NewCardFromString(`7c`),
		expSayGo:      false,
	}, {
		msg:           `when can't, still returns a card`,
		inputHand:     []string{`10c`},
		inputPrevPegs: []string{`10s`, `3c`, `4d`, `5s`, `6h`, `3h`, `2h`},
		inputCurPeg:   2,
		expCard:       cards.NewCardFromString(`10c`),
		expSayGo:      false,
	}, {
		msg:           `distinguishes between 10s`,
		inputHand:     []string{`10d`, `qh`, `ks`},
		inputPrevPegs: []string{`qs`, `9h`},
		inputCurPeg:   19,
		expCard:       cards.NewCardFromString(`10d`), // chosen as the first card in the slice
		expSayGo:      false,
	}, {
		msg:           `says go when can't play`,
		inputHand:     []string{`3s`, `2h`, `4s`},
		inputPrevPegs: []string{`10s`, `10h`, `10d`},
		inputCurPeg:   30,
		expCard:       cards.Card{},
		expSayGo:      true,
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
	}, {
		msg:           `when can't, still returns a card`,
		inputHand:     []string{`10c`},
		inputPrevPegs: []string{`10s`, `3c`, `4d`, `5s`, `6h`, `3h`, `2h`},
		inputCurPeg:   2,
		expCard:       cards.NewCardFromString(`10c`),
		expSayGo:      false,
	}, {
		msg:           `distinguishes between 10s`,
		inputHand:     []string{`10c`, `jd`, `qh`, `ks`},
		inputPrevPegs: []string{`qs`},
		inputCurPeg:   10,
		expCard:       cards.NewCardFromString(`qh`),
		expSayGo:      false,
	}, {
		msg:           `says go when can't play`,
		inputHand:     []string{`3s`, `2h`, `4s`},
		inputPrevPegs: []string{`10s`, `10h`, `10d`},
		inputCurPeg:   30,
		expCard:       cards.Card{},
		expSayGo:      true,
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
		inputHand:     []string{`2d`, `8h`, `jc`},
		inputPrevPegs: []string{`6d`, `7h`},
		inputCurPeg:   14,
		expCard:       cards.NewCardFromString(`8h`),
		expSayGo:      false,
	}, {
		msg:           `less obvious`,
		inputHand:     []string{`2d`, `7h`, `jc`},
		inputPrevPegs: []string{`6d`, `8h`},
		inputCurPeg:   14,
		expCard:       cards.NewCardFromString(`7h`),
		expSayGo:      false,
	}, {
		msg:           `least obvious (chooses the longer run)`,
		inputHand:     []string{`1d`, `4s`},
		inputPrevPegs: []string{`5d`, `6h`, `3s`, `2d`},
		inputCurPeg:   16,
		expCard:       cards.NewCardFromString(`4s`),
		expSayGo:      false,
	}, {
		msg:           `says go when can't play`,
		inputHand:     []string{`3s`, `2h`, `4s`},
		inputPrevPegs: []string{`10s`, `10h`, `10d`},
		inputCurPeg:   30,
		expCard:       cards.Card{},
		expSayGo:      true,
	}, {
		msg:           `when can't, still returns a card`,
		inputHand:     []string{`10c`},
		inputPrevPegs: []string{`10s`, `3c`, `4d`, `5s`, `6h`, `3h`, `2h`},
		inputCurPeg:   2,
		expCard:       cards.NewCardFromString(`10c`),
		expSayGo:      false,
	}}

	for _, tc := range testCases {
		actCard, actSayGo := PegToRun(strToCards(tc.inputHand), strToCards(tc.inputPrevPegs), tc.inputCurPeg)
		assert.Equal(t, tc.expSayGo, actSayGo, tc.msg)
		assert.Equal(t, tc.expCard, actCard, tc.msg)
	}
}
