package scorer

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshprzybyszewski/cribbage/model"
)

func BenchmarkHandPoints(b *testing.B) {
	hand := randomHand(b, 5)
	b.Run(`scoring random hand`, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s := HandPoints(hand[0], hand[1:])
			require.Less(b, s, 30)
			require.GreaterOrEqual(b, s, 0)
		}
	})
}

func TestResolveScoreType(t *testing.T) {
	testHands := scoringTestCases(t)
	for _, h := range testHands {
		h := h
		t.Run(`v1: `+h.desc, func(t *testing.T) {
			_, pts := scoreRunsAndPairs(h.values[:])
			assert.Equal(t, h.pointsByType[pointsTypeRuns]+h.pointsByType[pointsTypePairs], pts)
		})
		t.Run(`v2: `+h.desc, func(t *testing.T) {
			_, pts := scoreRunsAndPairsV2(h.values)
			assert.Equal(t, h.pointsByType[pointsTypeRuns]+h.pointsByType[pointsTypePairs], pts)
		})
	}
}

func TestPoints(t *testing.T) {
	testHands := scoringTestCases(t)
	for _, h := range testHands {
		h := h
		t.Run(h.desc, func(t *testing.T) {
			actPoints := HandPoints(h.cut, h.hand)
			assert.Equal(t, h.points, actPoints)
			desc, _ := pointsWithDesc(h.cut, h.hand, false)
			assert.Equal(t, h.scoreType, desc)
		})
	}
}

func TestScoringPoorlySizedHands(t *testing.T) {
	// Asserting zero also checks that the func doesn't panic
	assert.Zero(t, CribPoints(model.Card{}, make([]model.Card, 5)))
	assert.Zero(t, CribPoints(model.Card{}, make([]model.Card, 6)))
	assert.Zero(t, HandPoints(model.Card{}, make([]model.Card, 5)))
	assert.Zero(t, HandPoints(model.Card{}, make([]model.Card, 6)))
}

func randomHand(t testing.TB, n int) []model.Card {
	if n > 6 {
		t.Fatal(`you really don't need a hand with more than 6 cards, trust me`)
	}
	hand := make([]model.Card, 0, n)
	seen := make(map[model.Card]struct{}, n)
	for len(hand) < n {
		c := model.NewCardFromNumber(rand.Intn(52))
		if _, ok := seen[c]; ok {
			continue
		}
		hand = append(hand, c)
		seen[c] = struct{}{}
	}
	return hand
}
