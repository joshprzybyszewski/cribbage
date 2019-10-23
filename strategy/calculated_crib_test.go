package strategy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGiveCribHighestPotential(t *testing.T) {
	t.Skip(`haven't implemented tests correctly`)
	testCases := []struct {
		msg          string
		inputDesired int
		inputHand    []string
		expHand      []string
	}{{
		msg:          `obvious case`,
		inputDesired: 2,
		inputHand:    []string{`5s`, `5c`, `5d`, `5h`, `2h`, `1s`},
		expHand:      []string{`2h`, `1s`},
	}, {
		msg:          `obvious case still passses when requesting one card`,
		inputDesired: 1,
		inputHand:    []string{`5s`, `5c`, `5d`, `5h`, `1s`},
		expHand:      []string{`1s`},
	}, {
		msg:          `lots of tens`,
		inputDesired: 2,
		inputHand:    []string{`5s`, `JS`, `10s`, `10c`, `10d`, `10h`},
		expHand:      []string{`5s`},
	}, {
		msg:          `lots of 8s`,
		inputDesired: 2,
		inputHand:    []string{`5s`, `6s`, `7s`, `8s`, `8c`, `8d`},
		expHand:      []string{},
	}, {
		msg:          `only 7s and 8s`,
		inputDesired: 2,
		inputHand:    []string{`7s`, `7c`, `7d`, `8s`, `8c`, `8d`},
		expHand:      []string{},
	}, {
		msg:          `lots of 5s and 10s`,
		inputDesired: 2,
		inputHand:    []string{`5s`, `5c`, `5d`, `5h`, `10s`, `10c`},
		expHand:      []string{},
	}, {
		msg:          `lots of 5s with 7 and 8`,
		inputDesired: 2,
		inputHand:    []string{`5s`, `5c`, `5d`, `5h`, `7s`, `8c`},
		expHand:      []string{},
	}}

	for _, tc := range testCases {
		actHand := GiveCribHighestPotential(tc.inputDesired, strToCards(tc.inputHand))
		for _, c := range actHand {
			assert.True(t, containsCard(tc.inputHand, c), tc.msg+`: unexpected card `+c.String())
		}
		assert.Equal(t, strToCards(tc.expHand), actHand)
	}
}

func TestGiveCribLowestPotential(t *testing.T) {
	t.Skip(`haven't implemented tests correctly`)
	testCases := []struct {
		msg          string
		inputDesired int
		inputHand    []string
		expHand      []string
	}{{
		msg:          `obvious case`,
		inputDesired: 2,
		inputHand:    []string{`5s`, `10c`, `8d`, `9h`, `2h`, `1s`},
		expHand:      []string{},
	}, {
		msg:          `obvious case doesn't work when requesting one card`,
		inputDesired: 1,
		inputHand:    []string{`5s`, `10c`, `8d`, `9h`, `2h`, `1s`},
		expHand:      []string{},
	}, {
		msg:          `lots of tens`,
		inputDesired: 2,
		inputHand:    []string{`5s`, `JS`, `10s`, `10c`, `10d`, `10h`},
		expHand:      []string{},
	}, {
		msg:          `lots of 8s`,
		inputDesired: 2,
		inputHand:    []string{`5s`, `6s`, `7s`, `8s`, `8c`, `8d`},
		expHand:      []string{},
	}, {
		msg:          `only 7s and 8s`,
		inputDesired: 2,
		inputHand:    []string{`7s`, `7c`, `7d`, `8s`, `8c`, `8d`},
		expHand:      []string{},
	}, {
		msg:          `lots of 5s and 10s`,
		inputDesired: 2,
		inputHand:    []string{`5s`, `5c`, `5d`, `5h`, `10s`, `10c`},
		expHand:      []string{},
	}, {
		msg:          `lots of 5s with 7 and 8`,
		inputDesired: 2,
		inputHand:    []string{`5s`, `5c`, `5d`, `5h`, `7s`, `8c`},
		expHand:      []string{},
	}, {
		msg:          `has no 15 pair -- too low`,
		inputDesired: 2,
		inputHand:    []string{`5s`, `5c`, `5d`, `5h`, `6s`, `8c`},
		expHand:      []string{},
	}, {
		msg:          `has no 15 pair -- too high`,
		inputDesired: 2,
		inputHand:    []string{`10s`, `10c`, `10d`, `10h`, `Js`, `Qc`},
		expHand:      []string{},
	}, {
		msg:          `has no 15 pair -- way too high`,
		inputDesired: 2,
		inputHand:    []string{`Ks`, `Kc`, `Kd`, `Kh`, `Qs`, `Qc`},
		expHand:      []string{},
	}}

	for _, tc := range testCases {
		actHand := GiveCribLowestPotential(tc.inputDesired, strToCards(tc.inputHand))
		for _, c := range actHand {
			assert.True(t, containsCard(tc.inputHand, c), tc.msg+`: unexpected card `+c.String())
		}
		assert.Equal(t, strToCards(tc.expHand), actHand)
	}
}
