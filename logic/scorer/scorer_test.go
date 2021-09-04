package scorer

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joshprzybyszewski/cribbage/model"
)

func BenchmarkHandPoints(b *testing.B) {
	hand := make([]model.Card, 0, 5)
	seen := make(map[model.Card]struct{}, 5)
	for len(hand) < 5 {
		c := model.NewCardFromNumber(rand.Intn(52))
		if _, ok := seen[c]; ok {
			continue
		}
		hand = append(hand, c)
		seen[c] = struct{}{}
	}

	b.Run(`scoring random hand`, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s := HandPoints(hand[0], hand[1:])
			require.Less(b, s, 30)
			require.GreaterOrEqual(b, s, 0)
		}
	})
}

func TestScoreRuns(t *testing.T) {
	tests := []struct {
		desc      string
		hand      string
		expType   scoreType
		expPoints int
	}{{
		desc:      `no runs`,
		hand:      `5S,5C,5D,JH,5H`,
		expType:   none,
		expPoints: 0,
	}, {
		desc:      `triple run of 3`,
		hand:      `9S,8C,10D,8H,8H`,
		expType:   tripleRunOfThree,
		expPoints: 9,
	}, {
		desc:      `double run of 4`,
		hand:      `8S,8C,9D,10H,JH`,
		expType:   doubleRunOfFour,
		expPoints: 8,
	}, {
		desc:      `double double run of 3`,
		hand:      `8S,8C,9D,10H,9H`,
		expType:   doubleDoubleRunOfThree,
		expPoints: 12,
	}, {
		desc:      `double run of 3`,
		hand:      `8S,8C,9D,10H,KH`,
		expType:   doubleRunOfThree,
		expPoints: 6,
	}, {
		desc:      `run of 3`,
		hand:      `8S,2C,9D,10H,KH`,
		expType:   run3,
		expPoints: 3,
	}, {
		desc:      `run of 4`,
		hand:      `8S,JC,9D,10H,KH`,
		expType:   run4,
		expPoints: 4,
	}, {
		desc:      `run of 5`,
		hand:      `8S,JC,9D,10H,QH`,
		expType:   run5,
		expPoints: 5,
	}, {
		desc:      `random hand I got while playing`,
		hand:      `5H,3D,7D,7S,4C`,
		expType:   run3,
		expPoints: 3,
	}, {
		desc:      `just looking for ways to break it`,
		hand:      `1H,5D,7D,7S,9C`,
		expType:   none,
		expPoints: 0,
	}, {
		desc:      `another hand to break it`,
		hand:      `6D,6S,10H,9C,7H`,
		expType:   none,
		expPoints: 0,
	}, {
		desc:      `actual double run of three with a fifteen`,
		hand:      `6D,6S,10H,8C,7H`,
		expType:   doubleRunOfThree,
		expPoints: 6,
	}}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			h := parseHand(t, tc.hand)
			var vals [5]int
			var valuesToCounts valueToCount
			for i, c := range h {
				vals[i] = c.Value
				valuesToCounts[c.Value]++
			}
			st, pts := scoreRuns(vals, valuesToCounts)
			assert.Equal(t, tc.expType, st)
			assert.Equal(t, tc.expPoints, pts)
		})
	}
}

func TestScorePairs(t *testing.T) {
	tests := []struct {
		desc      string
		hand      string
		expPoints int
		expType   scoreType
	}{{
		desc:      `none`,
		hand:      `AC,2C,3C,4C,5C`,
		expPoints: 0,
		expType:   none,
	}, {
		desc:      `none, unsorted`,
		hand:      `AC,2C,3C,4C,5C`,
		expPoints: 0,
		expType:   none,
	}, {
		desc:      `quad`,
		hand:      `AC,2C,AH,AS,AD`,
		expPoints: 12,
		expType:   quad,
	}, {
		desc:      `triplet`,
		hand:      `AC,2C,AH,10S,AD`,
		expPoints: 6,
		expType:   triplet,
	}, {
		desc:      `two pair`,
		hand:      `AC,2C,AH,10S,10D`,
		expPoints: 4,
		expType:   twopair,
	}, {
		desc:      `one pair`,
		hand:      `3C,2C,AH,10S,10D`,
		expPoints: 2,
		expType:   onepair,
	}}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			h := parseHand(t, tc.hand)
			var vals [5]int
			var valuesToCounts valueToCount
			for i, c := range h {
				vals[i] = c.Value
				valuesToCounts[c.Value]++
			}
			st, pts := scorePairs(vals, valuesToCounts)
			assert.Equal(t, tc.expType, st)
			assert.Equal(t, tc.expPoints, pts)
		})
	}
}

func TestPointsStandardFunThings(t *testing.T) {
	testCases := []struct {
		desc      string
		leadCard  string
		hand      string
		expPoints int
	}{{
		desc:      `highest scoring hand`,
		leadCard:  `5H`,
		hand:      `5S,5C,5D,JH`,
		expPoints: 29,
	}, {
		desc:      `quad`,
		leadCard:  `8H`,
		hand:      `8S,8C,8D,10H`,
		expPoints: 12,
	}, {
		desc:      `triplet`,
		leadCard:  `KH`,
		hand:      `8S,8C,8D,10H`,
		expPoints: 6,
	}, {
		desc:      `one pair`,
		leadCard:  `KH`,
		hand:      `8S,8C,2D,10H`,
		expPoints: 2,
	}, {
		desc:      `two pair`,
		leadCard:  `KH`,
		hand:      `KS,8C,2D,8H`,
		expPoints: 4,
	}, {
		desc:      `triple run of 3`,
		leadCard:  `8H`,
		hand:      `9S,8C,10D,8H`,
		expPoints: 15,
	}, {
		desc:      `double run of 4`,
		leadCard:  `JH`,
		hand:      `8S,8C,9D,10H`,
		expPoints: 10,
	}, {
		desc:      `double double run of 3`,
		leadCard:  `9H`,
		hand:      `8S,8C,9D,10H`,
		expPoints: 16,
	}, {
		desc:      `double run of 3`,
		leadCard:  `KH`,
		hand:      `8S,8C,9D,10H`,
		expPoints: 8,
	}, {
		desc:      `run of 3`,
		leadCard:  `KH`,
		hand:      `8S,2C,9D,10H`,
		expPoints: 3,
	}, {
		desc:      `run of 4`,
		leadCard:  `KH`,
		hand:      `8S,JC,9D,10H`,
		expPoints: 4,
	}, {
		desc:      `run of 5`,
		leadCard:  `QH`,
		hand:      `8S,JC,9D,10H`,
		expPoints: 5,
	}, {
		desc:      `only nobs`,
		leadCard:  `6H`,
		hand:      `JH,KC,10D,8H`,
		expPoints: 1,
	}, {
		desc:      `flush`,
		leadCard:  `3H`,
		hand:      `8D,4D,10D,6D`,
		expPoints: 6,
	}, {
		desc:      `random hand I got while playing`,
		leadCard:  `4C`,
		hand:      `5H,3D,7D,7S`,
		expPoints: 9,
	}, {
		desc:      `just looking for ways to break it`,
		leadCard:  `9C`,
		hand:      `1H,5D,7D,7S`,
		expPoints: 6,
	}, {
		desc:      `another hand to break it`,
		leadCard:  `7H`,
		hand:      `6D,6S,10H,9C`,
		expPoints: 6,
	}, {
		desc:      `actual double run of three with a fifteen`,
		leadCard:  `7H`,
		hand:      `6D,6S,10H,8C`,
		expPoints: 10,
	}}

	for _, tc := range testCases {
		lead := model.NewCardFromString(tc.leadCard)
		hand := parseHand(t, tc.hand)

		actPoints := HandPoints(lead, hand)
		assert.Equal(t, tc.expPoints, actPoints, tc.desc)
	}
}

func TestPointsForFifteens(t *testing.T) {
	testCases := []struct {
		desc      string
		leadCard  string
		hand      string
		expPoints int
	}{{
		desc:      `highest scoring hand`,
		leadCard:  `5H`,
		hand:      `5S,5C,5D,JH`,
		expPoints: 29,
	}, {
		desc:      `15 for 2`,
		leadCard:  `8H`,
		hand:      `7S,AC,2D,KH`,
		expPoints: 2,
	}, {
		desc:      `run of 5 that adds up to 15`,
		leadCard:  `5H`,
		hand:      `AS,2C,3D,4H`,
		expPoints: 7,
	}, {
		desc:      `cards that up to under 15`,
		leadCard:  `4H`,
		hand:      `AS,2C,3D,4H`,
		expPoints: 10,
	}, {
		desc:      `cards that up to over 46`,
		leadCard:  `6H`,
		hand:      `KS,10C,QD,KH`,
		expPoints: 2,
	}}

	for _, tc := range testCases {
		lead := model.NewCardFromString(tc.leadCard)
		hand := parseHand(t, tc.hand)

		actPoints := HandPoints(lead, hand)
		assert.Equal(t, tc.expPoints, actPoints, tc.desc)
	}
}

func TestPointsOddInteractions(t *testing.T) {
	testCases := []struct {
		desc      string
		leadCard  string
		hand      string
		expPoints int
	}{{
		desc:      `flush, double run of 4, and fifteens`,
		leadCard:  `7D`,
		hand:      `6S,7S,8S,9S`,
		expPoints: 20,
	}, {
		desc:      `triplet across lead`,
		leadCard:  `8H`,
		hand:      `8S,8C,QD,10H`,
		expPoints: 6,
	}, {
		desc:      `run of 5 and flush`,
		leadCard:  `AS`,
		hand:      `3S,2S,5S,4S`,
		expPoints: 12,
	}}

	for _, tc := range testCases {
		lead := model.NewCardFromString(tc.leadCard)
		hand := make([]model.Card, 4)
		cardStrs := strings.Split(tc.hand, `,`)
		require.Len(t, cardStrs, 4)
		for i, c := range cardStrs {
			hand[i] = model.NewCardFromString(c)
		}

		actPoints := HandPoints(lead, hand)
		assert.Equal(t, tc.expPoints, actPoints, tc.desc)
	}
}

func TestScoringPoorlySizedHands(t *testing.T) {
	// Asserting zero also checks that the func doesn't panic
	assert.Zero(t, CribPoints(model.Card{}, make([]model.Card, 5)))
	assert.Zero(t, CribPoints(model.Card{}, make([]model.Card, 6)))
	assert.Zero(t, HandPoints(model.Card{}, make([]model.Card, 5)))
	assert.Zero(t, HandPoints(model.Card{}, make([]model.Card, 6)))
}

func parseHand(t *testing.T, handStr string) []model.Card {
	strs := strings.Split(handStr, `,`)
	hand := make([]model.Card, len(strs))
	for i, s := range strs {
		c, err := model.NewCardFromExternalString(strings.TrimSpace(s))
		require.NoError(t, err)
		hand[i] = c
	}
	return hand
}
