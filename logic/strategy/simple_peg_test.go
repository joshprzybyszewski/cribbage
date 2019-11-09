package strategy

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/joshprzybyszewski/cribbage/model"
)

func TestPegToFifteen(t *testing.T) {
	testCases := []struct {
		msg           string
		inputHand     []string
		inputPrevPegs []string
		inputCurPeg   int
		expCard       model.Card
		expSayGo      bool
	}{{
		msg:           `obvious`,
		inputHand:     []string{`3c`, `4d`, `5s`, `6h`},
		inputPrevPegs: []string{`10s`},
		inputCurPeg:   10,
		expCard:       model.NewCardFromString(`5s`),
		expSayGo:      false,
	}, {
		msg:           `when can't, still returns a card`,
		inputHand:     []string{`10c`},
		inputPrevPegs: []string{`10s`, `3c`, `4d`, `5s`, `6h`, `3h`, `2h`},
		inputCurPeg:   2,
		expCard:       model.NewCardFromString(`10c`),
		expSayGo:      false,
	}, {
		msg:           `when can't, still returns a card that is valid`,
		inputHand:     []string{`kh`, `1s`, `1d`},
		inputPrevPegs: []string{`10s`, `10c`, `9h`},
		inputCurPeg:   29,
		expCard:       model.NewCardFromString(`1s`), // the first valid card it can play
		expSayGo:      false,
	}, {
		msg:           `says go when can't play`,
		inputHand:     []string{`3s`, `2h`, `4s`},
		inputPrevPegs: []string{`10s`, `10h`, `10d`},
		inputCurPeg:   30,
		expCard:       model.Card{},
		expSayGo:      true,
	}}

	for _, tc := range testCases {
		actCard, actSayGo := PegToFifteen(strToCards(tc.inputHand), strToPeggedCards(tc.inputPrevPegs), tc.inputCurPeg)
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
		expCard       model.Card
		expSayGo      bool
	}{{
		msg:           `obvious`,
		inputHand:     []string{`ac`, `2h`, `4s`},
		inputPrevPegs: []string{`10s`, `10h`, `10d`},
		inputCurPeg:   30,
		expCard:       model.NewCardFromString(`ac`),
		expSayGo:      false,
	}, {
		msg:           `another less obvious`,
		inputHand:     []string{`7c`, `2h`, `4s`},
		inputPrevPegs: []string{`10s`, `10h`, `4d`},
		inputCurPeg:   24,
		expCard:       model.NewCardFromString(`7c`),
		expSayGo:      false,
	}, {
		msg:           `when can't, still returns a card`,
		inputHand:     []string{`10c`},
		inputPrevPegs: []string{`10s`, `3c`, `4d`, `5s`, `6h`, `3h`, `2h`},
		inputCurPeg:   2,
		expCard:       model.NewCardFromString(`10c`),
		expSayGo:      false,
	}, {
		msg:           `when can't, still returns a card that is valid`,
		inputHand:     []string{`kh`, `1s`, `1d`},
		inputPrevPegs: []string{`10s`, `10c`, `9h`},
		inputCurPeg:   29,
		expCard:       model.NewCardFromString(`1s`), // the first valid card it can play
		expSayGo:      false,
	}, {
		msg:           `distinguishes between 10s`,
		inputHand:     []string{`10d`, `qh`, `ks`},
		inputPrevPegs: []string{`qs`, `9h`},
		inputCurPeg:   19,
		expCard:       model.NewCardFromString(`10d`), // chosen as the first card in the slice
		expSayGo:      false,
	}, {
		msg:           `says go when can't play`,
		inputHand:     []string{`3s`, `2h`, `4s`},
		inputPrevPegs: []string{`10s`, `10h`, `10d`},
		inputCurPeg:   30,
		expCard:       model.Card{},
		expSayGo:      true,
	}}

	for _, tc := range testCases {
		actCard, actSayGo := PegToThirtyOne(strToCards(tc.inputHand), strToPeggedCards(tc.inputPrevPegs), tc.inputCurPeg)
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
		expCard       model.Card
		expSayGo      bool
	}{{
		msg:           `obvious`,
		inputHand:     []string{`3c`, `10s`, `6d`, `2h`},
		inputPrevPegs: []string{`3d`},
		inputCurPeg:   3,
		expCard:       model.NewCardFromString(`3c`),
		expSayGo:      false,
	}, {
		msg:           `when can't, still returns a card`,
		inputHand:     []string{`10c`},
		inputPrevPegs: []string{`10s`, `3c`, `4d`, `5s`, `6h`, `3h`, `2h`},
		inputCurPeg:   2,
		expCard:       model.NewCardFromString(`10c`),
		expSayGo:      false,
	}, {
		msg:           `when can't, still returns a valid card`,
		inputHand:     []string{`4c`, `7s`, `as`, `9c`},
		inputPrevPegs: []string{`10c`, `10s`, `8d`},
		inputCurPeg:   28,
		expCard:       model.NewCardFromString(`as`),
		expSayGo:      false,
	}, {
		msg:           `distinguishes between 10s`,
		inputHand:     []string{`10c`, `jd`, `qh`, `ks`},
		inputPrevPegs: []string{`qs`},
		inputCurPeg:   10,
		expCard:       model.NewCardFromString(`qh`),
		expSayGo:      false,
	}, {
		msg:           `says go when can't play`,
		inputHand:     []string{`3s`, `2h`, `4s`},
		inputPrevPegs: []string{`10s`, `10h`, `10d`},
		inputCurPeg:   30,
		expCard:       model.Card{},
		expSayGo:      true,
	}, {
		msg:           `still plays even if a pair seems possible, but is blocked by 31`,
		inputHand:     []string{`8s`, `10c`, `1s`},
		inputPrevPegs: []string{`9d`, `9h`, `10d`},
		inputCurPeg:   28,
		expCard:       model.NewCardFromString(`1s`),
		expSayGo:      false,
	}, {
		msg:           `says go when can't play even if a pair seems possible`,
		inputHand:     []string{`8s`, `10c`, `9s`},
		inputPrevPegs: []string{`9d`, `9h`, `10d`},
		inputCurPeg:   28,
		expCard:       model.Card{},
		expSayGo:      true,
	}}

	for _, tc := range testCases {
		actCard, actSayGo := PegToPair(strToCards(tc.inputHand), strToPeggedCards(tc.inputPrevPegs), tc.inputCurPeg)
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
		expCard       model.Card
		expSayGo      bool
	}{{
		msg:           `obvious`,
		inputHand:     []string{`2d`, `8h`, `jc`},
		inputPrevPegs: []string{`6d`, `7h`},
		inputCurPeg:   14,
		expCard:       model.NewCardFromString(`8h`),
		expSayGo:      false,
	}, {
		msg:           `less obvious`,
		inputHand:     []string{`2d`, `7h`, `jc`},
		inputPrevPegs: []string{`6d`, `8h`},
		inputCurPeg:   14,
		expCard:       model.NewCardFromString(`7h`),
		expSayGo:      false,
	}, {
		msg:           `least obvious (chooses the longer run)`,
		inputHand:     []string{`1d`, `4s`},
		inputPrevPegs: []string{`5d`, `6h`, `3s`, `2d`},
		inputCurPeg:   16,
		expCard:       model.NewCardFromString(`4s`),
		expSayGo:      false,
	}, {
		msg:           `doesn't try making a run over a go`,
		inputHand:     []string{`10c`, `5h`},
		inputPrevPegs: []string{`10d`, `10h`, `7s`, `6d`},
		inputCurPeg:   6,
		expCard:       model.NewCardFromString(`10c`), // first card in hand
		expSayGo:      false,
	}, {
		msg:           `says go when can't play`,
		inputHand:     []string{`3s`, `2h`, `4s`},
		inputPrevPegs: []string{`10s`, `10h`, `10d`},
		inputCurPeg:   30,
		expCard:       model.Card{},
		expSayGo:      true,
	}, {
		msg:           `when can't, still returns a card`,
		inputHand:     []string{`10c`},
		inputPrevPegs: []string{`10s`, `3c`, `4d`, `5s`, `6h`, `3h`, `2h`},
		inputCurPeg:   2,
		expCard:       model.NewCardFromString(`10c`),
		expSayGo:      false,
	}}

	for _, tc := range testCases {
		actCard, actSayGo := PegToRun(strToCards(tc.inputHand), strToPeggedCards(tc.inputPrevPegs), tc.inputCurPeg)
		assert.Equal(t, tc.expSayGo, actSayGo, tc.msg)
		assert.Equal(t, tc.expCard, actCard, tc.msg)
	}
}
