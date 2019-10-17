package strategy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAvoidCribFifteens(t *testing.T) {
	testCases := []struct {
		msg          string
		inputDesired int
		inputHand    []string
		canAvoid     bool
	}{{
		msg:          `obvious case`,
		inputDesired: 2,
		inputHand:    []string{`5s`, `5c`, `5d`, `5h`, `2h`, `1s`},
		canAvoid:     true,
	}, {
		msg:          `obvious case still passses when requesting one card`,
		inputDesired: 1,
		inputHand:    []string{`5s`, `5c`, `5d`, `5h`, `2h`, `1s`},
		canAvoid:     true,
	}, {
		msg:          `lots of tens`,
		inputDesired: 2,
		inputHand:    []string{`5s`, `JS`, `10s`, `10c`, `10d`, `10h`},
		canAvoid:     true,
	}, {
		msg:          `lots of 8s`,
		inputDesired: 2,
		inputHand:    []string{`5s`, `6s`, `7s`, `8s`, `8c`, `8d`},
		canAvoid:     true,
	}, {
		msg:          `only 7s and 8s`,
		inputDesired: 2,
		inputHand:    []string{`7s`, `7c`, `7d`, `8s`, `8c`, `8d`},
		canAvoid:     true,
	}, {
		msg:          `lots of 5s and 10s`,
		inputDesired: 2,
		inputHand:    []string{`5s`, `5c`, `5d`, `5h`, `10s`, `10c`},
		canAvoid:     true,
	}, {
		msg:          `lots of 5s with 7 and 8`,
		inputDesired: 2,
		inputHand:    []string{`5s`, `5c`, `5d`, `5h`, `7s`, `8c`},
		canAvoid:     true,
	}}

	for _, tc := range testCases {
		actHand := AvoidCribFifteens(tc.inputDesired, strToCards(tc.inputHand))
		sum := 0
		for _, c := range actHand {
			sum += c.PegValue()
		}
		if tc.canAvoid {
			assert.NotEqual(t, 15, sum)
		} else {
			assert.Equal(t, 15, sum)
		}
	}
}

func TestGiveCribFifteens(t *testing.T) {
	testCases := []struct {
		msg          string
		inputDesired int
		inputHand    []string
		canGive      bool
	}{{
		msg:          `obvious case`,
		inputDesired: 2,
		inputHand:    []string{`5s`, `10c`, `8d`, `9h`, `2h`, `1s`},
		canGive:      true,
	}, {
		msg:          `obvious case doesn't work when requesting one card`,
		inputDesired: 1,
		inputHand:    []string{`5s`, `10c`, `8d`, `9h`, `2h`, `1s`},
		canGive:      false,
	}, {
		msg:          `lots of tens`,
		inputDesired: 2,
		inputHand:    []string{`5s`, `JS`, `10s`, `10c`, `10d`, `10h`},
		canGive:      true,
	}, {
		msg:          `lots of 8s`,
		inputDesired: 2,
		inputHand:    []string{`5s`, `6s`, `7s`, `8s`, `8c`, `8d`},
		canGive:      true,
	}, {
		msg:          `only 7s and 8s`,
		inputDesired: 2,
		inputHand:    []string{`7s`, `7c`, `7d`, `8s`, `8c`, `8d`},
		canGive:      true,
	}, {
		msg:          `lots of 5s and 10s`,
		inputDesired: 2,
		inputHand:    []string{`5s`, `5c`, `5d`, `5h`, `10s`, `10c`},
		canGive:      true,
	}, {
		msg:          `lots of 5s with 7 and 8`,
		inputDesired: 2,
		inputHand:    []string{`5s`, `5c`, `5d`, `5h`, `7s`, `8c`},
		canGive:      true,
	}, {
		msg:          `has no 15 pair -- too low`,
		inputDesired: 2,
		inputHand:    []string{`5s`, `5c`, `5d`, `5h`, `6s`, `8c`},
		canGive:      false,
	}, {
		msg:          `has no 15 pair -- too high`,
		inputDesired: 2,
		inputHand:    []string{`10s`, `10c`, `10d`, `10h`, `Js`, `Qc`},
		canGive:      false,
	}, {
		msg:          `has no 15 pair -- way too high`,
		inputDesired: 2,
		inputHand:    []string{`Ks`, `Kc`, `Kd`, `Kh`, `Qs`, `Qc`},
		canGive:      false,
	}}

	for _, tc := range testCases {
		actHand := GiveCribFifteens(tc.inputDesired, strToCards(tc.inputHand))
		sum := 0
		for _, c := range actHand {
			sum += c.PegValue()
		}
		if tc.canGive {
			assert.Equal(t, 15, sum)
		} else {
			assert.NotEqual(t, 15, sum)
		}
	}
}

func TestAvoidCribPairs(t *testing.T) {
	testCases := []struct {
		msg       string
		inputHand []string
		canAvoid  bool
	}{{
		msg:       `obvious case`,
		inputHand: []string{`5s`, `5c`, `5d`, `5h`, `2h`, `1s`},
		canAvoid:  true,
	}, {
		msg:       `lots of tens`,
		inputHand: []string{`5s`, `JS`, `10s`, `10c`, `10d`, `10h`},
		canAvoid:  true,
	}, {
		msg:       `lots of 8s`,
		inputHand: []string{`5s`, `6s`, `7s`, `8s`, `8c`, `8d`},
		canAvoid:  true,
	}, {
		msg:       `only 7s and 8s`,
		inputHand: []string{`7s`, `7c`, `7d`, `8s`, `8c`, `8d`},
		canAvoid:  false,
	}, {
		msg:       `lots of 5s and 10s`,
		inputHand: []string{`5s`, `5c`, `5d`, `5h`, `10s`, `10c`},
		canAvoid:  false,
	}, {
		msg:       `lots of 5s with 7 and 8`,
		inputHand: []string{`5s`, `5c`, `5d`, `5h`, `7s`, `8c`},
		canAvoid:  true,
	}}

	for _, tc := range testCases {
		actHand := AvoidCribPairs(2, strToCards(tc.inputHand))
		if tc.canAvoid {
			assert.NotEqual(t, actHand[0].Value, actHand[1].Value)
		} else {
			assert.Equal(t, actHand[0].Value, actHand[1].Value)
		}
	}

	actHand := AvoidCribPairs(1, strToCards([]string{`7s`, `8c`, `9d`, `10h`, `Js`, `Qc`}))
	assert.Len(t, actHand, 1)
}

func TestGiveCribPairs(t *testing.T) {
	testCases := []struct {
		msg       string
		inputHand []string
		canGive   bool
	}{{
		msg:       `obvious case`,
		inputHand: []string{`5s`, `5c`, `8d`, `9h`, `2h`, `1s`},
		canGive:   true,
	}, {
		msg:       `lots of tens`,
		inputHand: []string{`5s`, `JS`, `10s`, `10c`, `10d`, `10h`},
		canGive:   true,
	}, {
		msg:       `lots of 8s`,
		inputHand: []string{`5s`, `6s`, `7s`, `8s`, `8c`, `8d`},
		canGive:   true,
	}, {
		msg:       `only 7s and 8s`,
		inputHand: []string{`7s`, `7c`, `7d`, `8s`, `8c`, `8d`},
		canGive:   true,
	}, {
		msg:       `lots of 5s and 10s`,
		inputHand: []string{`5s`, `5c`, `5d`, `5h`, `10s`, `10c`},
		canGive:   true,
	}, {
		msg:       `lots of 5s with 7 and 8`,
		inputHand: []string{`5s`, `5c`, `5d`, `5h`, `7s`, `8c`},
		canGive:   true,
	}, {
		msg:       `has no pair -- too low`,
		inputHand: []string{`5s`, `4c`, `3d`, `2h`, `6s`, `8c`},
		canGive:   false,
	}, {
		msg:       `has no pair -- too high`,
		inputHand: []string{`7s`, `8c`, `9d`, `10h`, `Js`, `Qc`},
		canGive:   false,
	}}

	for _, tc := range testCases {
		actHand := GiveCribPairs(2, strToCards(tc.inputHand))
		if tc.canGive {
			assert.Equal(t, actHand[0].Value, actHand[1].Value)
		} else {
			assert.NotEqual(t, actHand[0].Value, actHand[1].Value)
		}
	}

	actHand := GiveCribPairs(1, strToCards([]string{`7s`, `8c`, `9d`, `10h`, `Js`, `Qc`}))
	assert.Len(t, actHand, 1)
}
