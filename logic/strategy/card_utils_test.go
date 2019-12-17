package strategy

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshprzybyszewski/cribbage/model"
)

func validateHand(superset, subset []model.Card) bool {
	thisHandMap := make(map[model.Card]struct{}, len(subset))
	for _, c := range subset {
		_, ok := thisHandMap[c]
		if ok {
			return false
		}
		thisHandMap[c] = struct{}{}
	}
	ct := 0
	for _, c := range superset {
		_, ok := thisHandMap[c]
		if ok {
			ct++
		}
	}
	return ct == len(subset)
}

func factorial(n int, cache map[int]int) int {
	res, ok := cache[n]
	if ok {
		return res
	}
	if n == 1 {
		res = 1
	} else {
		res = n * factorial(n-1, cache)
	}
	cache[n] = res
	return res
}

func nchoosek(n, k int, fCache map[int]int) (int, error) {
	if k > n {
		return 0, errors.New(`k must be less than or equal to n`)
	}
	return factorial(n, fCache) / (factorial(k, fCache) * factorial(n-k, fCache)), nil
}

func generateHand(n int) []model.Card {
	hand := make([]model.Card, n)
	for i := 0; i < n; i++ {
		hand[i] = model.NewCardFromNumber(i)
	}
	return hand
}

func TestChooseFrom(t *testing.T) {
	tests := []struct {
		desc   string
		hand   []model.Card
		nCards int
		expErr string
	}{{
		desc:   `6 choose 2`,
		hand:   generateHand(6),
		nCards: 2,
		expErr: ``,
	}, {
		desc:   `6 choose 3`,
		hand:   generateHand(6),
		nCards: 3,
		expErr: ``,
	}, {
		desc:   `6 choose 4`,
		hand:   generateHand(6),
		nCards: 4,
		expErr: ``,
	}, {
		desc:   `5 choose 2`,
		hand:   generateHand(5),
		nCards: 2,
		expErr: ``,
	}, {
		desc:   `5 choose 3`,
		hand:   generateHand(5),
		nCards: 3,
		expErr: ``,
	}, {
		desc:   `5 choose 4`,
		hand:   generateHand(5),
		nCards: 4,
		expErr: ``,
	}, {
		desc:   `5 choose 6`,
		hand:   generateHand(5),
		nCards: 6,
		expErr: `developer error: invalid k`,
	}, {
		desc:   `choose zero cards`,
		hand:   generateHand(5),
		nCards: 0,
		expErr: `developer error: invalid k`,
	}, {
		desc:   `hand too large`,
		hand:   generateHand(7),
		nCards: 3,
		expErr: `too many cards in hand (maximum 6)`,
	}}
	fCache := make(map[int]int)
	for _, tc := range tests {
		all, err := chooseFrom(tc.nCards, tc.hand)
		if tc.expErr != `` {
			assert.EqualError(t, err, tc.expErr)
		} else {
			assert.Nil(t, err)
			expNum, err := nchoosek(len(tc.hand), tc.nCards, fCache)
			require.Nil(t, err)
			assert.Equal(t, expNum, len(all))
			for _, h := range all {
				assert.Equal(t, tc.nCards, len(h))
				ok := validateHand(tc.hand, h)
				assert.True(t, ok)
			}
		}
	}
}
